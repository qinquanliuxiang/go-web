package service

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"qqlx/base/apierr"
	"qqlx/base/conf"
	"qqlx/base/helpers"
	"qqlx/base/interfaces"
	"qqlx/base/logger"
	"qqlx/base/reason"
	"qqlx/model"
	"qqlx/pkg/jwt"
	"qqlx/pkg/sonyflake"
	"qqlx/schema"
	"qqlx/store/cache"
	"qqlx/store/rbac"
	"qqlx/store/userstore"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSVC struct {
	generateID    *sonyflake.GenerateIDStruct
	userStore     interfaces.UserStoreInterface
	userRoleStore interfaces.UserRoleStoreInterface
	roleStore     interfaces.RoleStoreInterface
	cache         interfaces.CacheInterface
	casbin        interfaces.CasbinInterface
	salt          string
	ldapEnable    bool
	ldap          interfaces.LdapInterface
}

func NewUserSVC(
	generateID *sonyflake.GenerateIDStruct, userStore interfaces.UserStoreInterface, userRoleStore interfaces.UserRoleStoreInterface, roleStore interfaces.RoleStoreInterface, cache interfaces.CacheInterface, casbin interfaces.CasbinInterface, ldap interfaces.LdapInterface) (*UserSVC, error) {
	ldapEnable := conf.GetLdapEnable()
	salt, err := conf.GetSalt()
	if err != nil {
		return nil, err
	}
	userSvc := &UserSVC{
		generateID:    generateID,
		userStore:     userStore,
		userRoleStore: userRoleStore,
		roleStore:     roleStore,
		cache:         cache,
		casbin:        casbin,
		salt:          salt,
		ldap:          ldap,
		ldapEnable:    ldapEnable,
	}
	return userSvc, nil
}

func (receive *UserSVC) RegistryUser(ctx context.Context, req *schema.UserRegistryRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("users registry request: %#v", req)
	var (
		user            *model.User
		encryptPassword string
		id              int
	)

	user, err = receive.userStore.Query(ctx, userstore.Email(req.Email))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 创建 ldap 用户
		if receive.ldapEnable {
			// 生成 ldap ssha 密码
			ssha := receive.ldapEncryptSSHA(req.Password)
			// 创建 ldap 用户
			err = receive.ldap.CreateUser(ctx, req.Name, ssha, req.Email)
			if err != nil {
				return err
			}
		}

		encryptPassword, err = receive.encryptPassword(ctx, req.Password)
		if err != nil {
			return err
		}
		id, err = receive.generateID.NextID()
		if err != nil {
			return err
		}
		err = receive.userStore.Create(ctx, &model.User{
			ID:       id,
			Name:     req.Name,
			NickName: req.NickName,
			Password: encryptPassword,
			Avatar:   req.Avatar,
			Email:    req.Email,
			Mobile:   req.Mobile,
		})
		if err != nil {
			return err
		}

	}
	if user != nil {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, "user already exists", reason.ErrUserExists)
	}
	return nil
}

func (receive *UserSVC) Login(ctx context.Context, req *schema.UserLoginRequest) (res *schema.UserLoginResponse, err error) {
	logger.WithContext(ctx, true).Debugf("user login, request: %#v", req)
	var user *model.User
	if req.Email != "" {
		user, err = receive.userStore.Query(ctx, userstore.Email(req.Email), userstore.LoadRoles())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apierr.Unauthorized().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserNotFound)
			}
			return nil, err
		}
	}

	if *user.Status == model.UserStatusDisable {
		logger.WithContext(ctx, true).Errorf("users has been disabled, user email: %s", user.Email)
		return nil, apierr.Unauthorized().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserIsDisable)
	}
	if !receive.verifyPassword(ctx, req.Password, user.Password) {
		return nil, apierr.Unauthorized().Set(apierr.ServiceErrCode, "invalid password", reason.ErrInvalidPassword)
	}
	roleCount := len(user.Roles)
	if roleCount > 0 {
		rolesName := make([]any, 0, roleCount)
		for _, role := range user.Roles {
			rolesName = append(rolesName, role.Name)
		}
		err = receive.cache.SetSet(ctx, helpers.GetRoleCacheKey(user.Name), rolesName, &cache.NeverExpires)
		if err != nil {
			return nil, err
		}
	}

	token, err := jwt.NewClaims(user.ID, user.Name).GenerateToken()
	if err != nil {
		return nil, err
	}
	res = &schema.UserLoginResponse{
		Token: token,
	}
	return res, err
}

func (receive *UserSVC) Logout(ctx context.Context, id int) (err error) {
	query, _ := receive.userStore.Query(ctx, userstore.ID(id))
	if query.Name != "" {
		_ = receive.cache.Del(ctx, helpers.GetRoleCacheKey(query.Name))
	}
	return nil
}

func (receive *UserSVC) DisableUser(ctx context.Context, req *schema.UserQueryRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("user delete, request: %#v", req)
	var user *model.User
	user, err = receive.userStore.Query(ctx, userstore.ID(req.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierr.InternalServer().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserNotFound)
		}
		return err
	}
	if user.Name == "admin" {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, reason.ErrAdminUserNotAllow.Error(), reason.ErrAdminUserNotAllow)
	}

	// 删除 ldap 用户
	if receive.ldapEnable {
		// 删除用户后，所在组中的记录也会被删除
		err = receive.ldap.DeleteUser(ctx, user.Name)
		if err != nil {
			return err
		}
	}

	if *user.Status == model.UserStatusDisable {
		logger.WithContext(ctx, true).Errorf("user has been disabled, userName: %s", user.Name)
		return apierr.InternalServer().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserIsDisable)
	}

	user.Status = &model.UserStatusDisable
	err = receive.userStore.Save(ctx, user)
	if err != nil {
		return err
	}

	return receive.cache.Del(ctx, helpers.GetRoleCacheKey(user.Name))
}

func (receive *UserSVC) EnableUser(ctx context.Context, req *schema.UserEnableRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("user enable, request: %#v", req)
	var user *model.User
	user, err = receive.userStore.Query(ctx, userstore.ID(req.ID), userstore.LoadRoles())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierr.InternalServer().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserNotFound)
		}
		return err
	}
	if *user.Status != model.UserStatusDisable {
		logger.WithContext(ctx, true).Errorf("user has been enabled, userName: %s", user.Name)
		return nil
	}

	// 添加 ldap 用户
	if receive.ldapEnable {
		ldapPassword := receive.ldapEncryptSSHA(req.Password)
		err = receive.ldap.CreateUser(ctx, user.Name, ldapPassword, user.Email)
		if err != nil {
			return err
		}

		roleCount := len(user.Roles)
		if roleCount > 0 {
			for _, role := range user.Roles {
				err = receive.ldap.AddUserToGroup(ctx, role.Name, user.Name)
				if err != nil {
					return err
				}
			}
		}
	}

	user.Status = &model.UserStatusAvailable
	password, err := receive.encryptPassword(ctx, req.Password)
	if err != nil {
		return err
	}

	user.Password = password
	return receive.userStore.Save(ctx, user)
}

func (receive *UserSVC) UpdatePassword(ctx context.Context, req *schema.UserUpdatePasswordRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("user update password, request: %#v", req)
	var user *model.User
	user, err = receive.userStore.Query(ctx, userstore.ID(req.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierr.InternalServer().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserNotFound)
		}
		return err
	}
	if *user.Status == model.UserStatusDisable {
		logger.WithContext(ctx, true).Errorf("user has been disabled, user email: %s", user.Email)
		return apierr.InternalServer().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserIsDisable)
	}

	if !receive.verifyPassword(ctx, req.OldPassword, user.Password) {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, "invalid password", reason.ErrInvalidPassword)
	}

	// 更新 ldap 用户
	if receive.ldapEnable {
		err = receive.ldap.UpdateUserPassword(ctx, user.Name, req.NewPassword)
		if err != nil {
			return err
		}
	}

	encryptPassword, err := receive.encryptPassword(ctx, req.NewPassword)
	if err != nil {
		return err
	}
	user.Password = encryptPassword

	return receive.userStore.Save(ctx, user)
}

func (receive *UserSVC) UpdateUser(ctx context.Context, req *schema.UserUpdateRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("user update, request: %#v", req)
	var user *model.User
	user, err = receive.userStore.Query(ctx, userstore.ID(req.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierr.InternalServer().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserNotFound)
		}
		return err
	}

	isUpdated := false
	if req.Mobile != "" {
		user.Mobile = req.Mobile
		isUpdated = true
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
		isUpdated = true
	}
	if req.NickName != "" {
		user.NickName = req.NickName
		isUpdated = true
	}

	if !isUpdated {
		return nil
	}
	return receive.userStore.Save(ctx, user)
}

// UserAddRole 增加用户角色
func (receive *UserSVC) UserAddRole(ctx context.Context, req *schema.UserUpdateRoleRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("user update role, request: %#v", req)
	roleNames := helpers.Deduplicate(req.RoleNames)
	var user *model.User
	user, err = receive.userStore.Query(ctx, userstore.ID(req.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierr.InternalServer().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserNotFound)
		}
		return err
	}
	if *user.Status == model.UserStatusDisable {
		logger.WithContext(ctx, true).Errorf("user has been disabled, userName: %s", user.Name)
		return apierr.InternalServer().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserIsDisable)
	}

	if user.Name == "admin" {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, reason.ErrAdminUserNotAllow.Error(), reason.ErrAdminUserNotAllow)
	}

	roleCount := len(roleNames)
	_, list, err := receive.roleStore.List(ctx, 1, roleCount, rbac.RoleNames(roleNames))
	if err != nil {
		return err
	}
	notFound := helpers.FindMissingByName(list, roleNames)
	if len(notFound) > 0 {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, fmt.Sprintf("role not exist: %v", notFound), reason.ErrRoleNotFound)
	}

	userCache := helpers.GetRoleCacheKey(user.Name)
	if err = receive.cache.Del(ctx, userCache); err != nil {
		return err
	}

	// 更新 ldap
	if receive.ldapEnable {
		for _, roleName := range roleNames {
			var exist bool
			exist, err = receive.ldap.SearchGroup(ctx, roleName)
			if err != nil {
				return err
			}
			if !exist {
				err = receive.ldap.CreateGroup(ctx, roleName)
				if err != nil {
					return err
				}
			}
			err = receive.ldap.AddUserToGroup(ctx, roleName, user.Name)
			if err != nil {
				return err
			}
		}
	}

	err = receive.userRoleStore.AppendRoles(ctx, user, list)
	if err != nil {
		return err
	}

	query, err := receive.userStore.Query(ctx, userstore.ID(req.ID), userstore.LoadRoles())
	if err != nil {
		return err
	}

	if len(query.Roles) > 0 {
		listNames := make([]any, 0, len(list))
		for _, role := range query.Roles {
			listNames = append(listNames, role.Name)
		}
		if err = receive.cache.SetSet(ctx, userCache, listNames, &cache.NeverExpires); err != nil {
			return err
		}
	}

	go func() {
		time.Sleep(time.Millisecond * 200)
		if err = receive.cache.Del(ctx, userCache); err != nil {
			logger.WithContext(ctx, true).Error(err)
		}
	}()

	return nil
}

func (receive *UserSVC) UserRemoveRole(ctx context.Context, req *schema.UserUpdateRoleRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("user remove role, request: %#v", req)
	uniqRoleNames := helpers.Deduplicate(req.RoleNames)
	var user *model.User
	user, err = receive.userStore.Query(ctx, userstore.ID(req.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierr.InternalServer().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserNotFound)
		}
		return err
	}

	if user.Name == "admin" {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, reason.ErrAdminUserNotAllow.Error(), reason.ErrAdminUserNotAllow)
	}

	if *user.Status == model.UserStatusDisable {
		logger.WithContext(ctx, true).Errorf("user has been disabled, userName: %s", user.Name)
		return apierr.InternalServer().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserIsDisable)
	}

	roleCount := len(uniqRoleNames)
	_, list, err := receive.roleStore.List(ctx, 1, roleCount, rbac.RoleNames(uniqRoleNames))
	if err != nil {
		return err
	}

	notFound := helpers.FindMissingByName(list, uniqRoleNames)
	if len(notFound) > 0 {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, fmt.Sprintf("role not exist: %v", notFound), reason.ErrRoleNotFound)
	}

	// 删除 ldap 用户
	if receive.ldapEnable {
		for _, roleName := range uniqRoleNames {
			err = receive.ldap.RemoveUserFromGroup(ctx, roleName, user.Name)
			if err != nil {
				return err
			}
		}
	}

	userCache := helpers.GetRoleCacheKey(user.Name)
	if err = receive.cache.Del(ctx, userCache); err != nil {
		return err
	}

	err = receive.userRoleStore.DeleteRoles(ctx, user, list)
	if err != nil {
		return err
	}

	query, err := receive.userStore.Query(ctx, userstore.ID(user.ID), userstore.LoadRoles())
	if err != nil {
		return err
	}
	queryRoleCount := len(query.Roles)
	if queryRoleCount > 0 {
		roleNames := make([]any, 0, queryRoleCount)
		for _, role := range query.Roles {
			roleNames = append(roleNames, role.Name)
		}
		if err = receive.cache.SetSet(ctx, userCache, roleNames, &cache.NeverExpires); err != nil {
			return err
		}
	}

	go func() {
		time.Sleep(time.Millisecond * 200)
		if err = receive.cache.Del(ctx, userCache); err != nil {
			logger.WithContext(ctx, true).Error(err)
		}
	}()

	return nil
}

// Info 获取用户信息
func (receive *UserSVC) Info(ctx context.Context, req *schema.UserQueryRequest) (res *schema.UserResponse, err error) {
	logger.WithContext(ctx, true).Debugf("user info, request: %#v", req)
	options := make([]userstore.QueryOption, 0, len(req.Query)+1)
	if len(req.Query) > 0 {
		for _, v := range req.Query {
			switch v {
			case "roles":
				options = append(options, userstore.LoadRoles())
			}
		}
	}
	options = append(options, userstore.ID(req.ID))
	user, err := receive.userStore.Query(ctx, options...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.InternalServer().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserNotFound)
		}
		return nil, err
	}
	if *user.Status == model.UserStatusDisable {
		logger.WithContext(ctx, true).Errorf("userstore has been disabled, userName: %s", user.Name)
		return nil, apierr.InternalServer().Set(apierr.ServiceErrCode, "user not found", reason.ErrUserIsDisable)
	}

	res = &schema.UserResponse{}
	res.ConvertToUserResponse(user)
	if len(req.Query) == 0 {
		roleName, err := receive.cache.GetSet(ctx, helpers.GetRoleCacheKey(user.Name))
		if err != nil {
			return nil, err
		}
		res.RoleName = roleName
	}
	return res, nil
}

func (receive *UserSVC) ListUser(ctx context.Context, req *schema.UserListRequest) (data *schema.UserListResponse, err error) {
	logger.WithContext(ctx, true).Debugf("user list, request: %#v", req)
	options := make([]userstore.QueryOption, 0)
	// 过滤关键字
	if req.Keyword != "" {
		options = append(options, userstore.QueryByNameOrEmail(req.Keyword, req.Value))
	}
	// 过滤状态
	options = append(options, userstore.Status(req.Status), userstore.SortByCreatedDesc())

	total, users, err := receive.userStore.List(ctx, req.Page, req.PageSize, options...)
	if err != nil {
		return nil, err
	}
	return receive.formatUserList(ctx, req, total, users), nil
}

func (receive *UserSVC) formatUserList(_ context.Context, req *schema.UserListRequest, total int64, users []model.User) *schema.UserListResponse {
	res := &schema.UserListResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		Items:    make([]schema.UserResponse, 0, len(users)),
	}

	for i := range users {
		userRes := schema.UserResponse{}
		userRes.ConvertToUserResponse(&users[i])
		res.Items = append(res.Items, userRes)
	}

	return res
}

// encryptPassword 加密密码
func (receive *UserSVC) encryptPassword(_ context.Context, Pass string) (string, error) {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(Pass), bcrypt.DefaultCost)
	if err != nil {
		logger.WithContext(context.Background(), true).Errorf("failed to encrypt password: %s, error: %v", Pass, err)
		return "", apierr.InternalServer().Set(apierr.ServiceErrCode, "unknown error", reason.ErrEncryptPassword)
	}
	return string(hashPwd), nil
}

// verifyPassword 验证密码
func (receive *UserSVC) verifyPassword(_ context.Context, loginPass, userPass string) bool {
	if len(loginPass) == 0 && len(userPass) == 0 {
		return true
	}
	err := bcrypt.CompareHashAndPassword([]byte(userPass), []byte(loginPass))
	return err == nil
}

func (receive *UserSVC) ldapEncryptSSHA(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	hash.Write([]byte(receive.salt))

	hashResult := hash.Sum(nil)
	result := append(hashResult, []byte(receive.salt)...)
	encoded := base64.StdEncoding.EncodeToString(result)
	return "{SSHA}" + encoded
}
