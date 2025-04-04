package service

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"qqlx/base/apierr"
	"qqlx/base/conf"
	"qqlx/base/constant"
	"qqlx/base/helpers"
	"qqlx/base/logger"
	"qqlx/base/reason"
	"qqlx/model"
	"qqlx/pkg/jwt"
	"qqlx/pkg/sonyflake"
	"qqlx/schema"
	"qqlx/store"
	"qqlx/store/cache"
	"qqlx/store/rbac"
	"qqlx/store/userstore"
	"time"
)

type UserSVC struct {
	generateID *sonyflake.GenerateIDStruct
	userStore  store.UserStoreInterface
	roleStore  store.RoleStoreInterface
	cache      store.CacheInterface
	casbin     store.CasbinInterface
	salt       string
	ldapEnable bool
	ldap       store.LdapInterface
}

func NewUserSVC(
	generateID *sonyflake.GenerateIDStruct, userStore store.UserStoreInterface, roleStore store.RoleStoreInterface, cache store.CacheInterface, casbin store.CasbinInterface, ldap store.LdapInterface) (*UserSVC, error) {
	ldapEnable := conf.GetLdapEnable()
	salt, err := conf.GetSalt()
	if err != nil {
		return nil, err
	}
	userSvc := &UserSVC{
		generateID: generateID,
		userStore:  userStore,
		roleStore:  roleStore,
		cache:      cache,
		casbin:     casbin,
		salt:       salt,
		ldap:       ldap,
		ldapEnable: ldapEnable,
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
		//var roleExist bool
		//role, err := receive.roleStore.Query(ctx, rbac.Name(req.Name))
		//if err != nil {
		//	if !errors.Is(err, gorm.ErrRecordNotFound) {
		//		return apierr.InternalServer().WithStack().WithErr(err)
		//	}
		//}
		//if role != nil {
		//	roleExist = true
		//}
		//if !roleExist {
		//	return apierr.InternalServer().WithMsg("role not exist").WithErr(err)
		//}

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
			//Name: req.Name,
			Password: encryptPassword,
			Avatar:   req.Avatar,
			Email:    req.Email,
			Mobile:   req.Mobile,
		})
		if err != nil {
			return err
		}

		// 生成 ldap ssha 密码
		ssha := receive.ldapEncryptSSHA(req.Password)
		if receive.ldapEnable {
			// 如果用户组不为空,添加用户到组
			//if req.Name != "" {
			//	exist, err := receive.ldap.SearchGroup(ctx, req.Name)
			//	if err != nil {
			//		return apierr.InternalServer().WithStack().WithMsg("create user failed").WithErr(err)
			//	}
			//	if !exist {
			//		err = receive.ldap.CreateGroup(ctx, req.Name)
			//		if err != nil {
			//			return apierr.InternalServer().WithStack().WithMsg("create user failed").WithErr(err)
			//		}
			//	}
			//	err = receive.ldap.AddUserToGroup(ctx, req.Name, req.Name)
			//	if err != nil {
			//		return apierr.InternalServer().WithStack().WithMsg("user add to group failed").WithErr(err)
			//	}
			//}
			// 创建 ldap 用户
			err = receive.ldap.CreateUser(ctx, req.Name, ssha, req.Email)
			if err != nil {
				return apierr.InternalServer().WithStack().WithMsg("create user failed").WithErr(err)
			}
		}
	}

	if user != nil {
		return apierr.InternalServer().WithErr(fmt.Errorf("user already exists")).WithStack()
	}
	return nil
}

func (receive *UserSVC) Login(ctx context.Context, req *schema.UserLoginRequest) (res *schema.UserLoginResponse, err error) {
	logger.WithContext(ctx, true).Debugf("user login, request: %#v", req)
	var user *model.User
	if req.Email != "" {
		user, err = receive.userStore.Query(ctx, userstore.Email(req.Email), userstore.LoadRole())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apierr.Unauthorized().WithStack().WithErr(reason.ErrUserNotFound)
			}
			return nil, err
		}
	}

	if user == nil {
		return nil, apierr.Unauthorized().WithStack().WithErr(reason.ErrUserNotFound)
	}

	if *user.Status == model.UserStatusDisable {
		logger.WithContext(ctx, true).Errorf("users has been disabled, user email: %s", user.Email)
		return nil, apierr.Unauthorized().WithStack().WithErr(reason.ErrUserNotFound)
	}
	if !receive.verifyPassword(ctx, req.Password, user.Password) {
		return nil, apierr.Unauthorized().WithStack().WithErr(reason.ErrInvalidPassword)
	}

	if user.RoleName != "" {
		err = receive.cache.SetString(ctx, constant.RoleCacheKeyPrefix+user.Name, user.Role.Name, &cache.NeverExpires)
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
	query, err := receive.userStore.Query(ctx, userstore.ID(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierr.InternalServer().WithStack().WithErr(reason.ErrUserNotFound)
		}
		return err
	}
	if query != nil {
		err = receive.cache.Del(ctx, helpers.GetRoleCacheKey(query.Name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (receive *UserSVC) DeleteUser(ctx context.Context, req *schema.UserIDRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("user delete, request: %#v", req)
	var user *model.User
	user, err = receive.userStore.Query(ctx, userstore.ID(req.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierr.InternalServer().WithStack().WithErr(reason.ErrUserNotFound)
		}
		return err
	}
	if *user.Status == model.UserStatusDisable {
		logger.WithContext(ctx, true).Errorf("user has been disabled, userName: %s", user.Name)
		return apierr.InternalServer().WithStack().WithErr(reason.ErrUserNotFound)
	}

	user.Status = &model.UserStatusDisable
	err = receive.userStore.Save(ctx, user)
	if err != nil {
		return err
	}

	err = receive.cache.Del(ctx, helpers.GetRoleCacheKey(user.Name))
	if err != nil {
		return err
	}

	if receive.ldapEnable {
		// 删除用户后，所在组中的记录也会被删除
		err = receive.ldap.DeleteUser(ctx, user.Name)
		if err != nil {
			return err
		}
		//err = receive.ldap.RemoveUserFromGroup(ctx, user.Name, user.Name)
		//if err != nil {
		//	return err
		//}
	}
	return nil
}

func (receive *UserSVC) UpdatePassword(ctx context.Context, req *schema.UserUpdatePasswordRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("user update password, request: %#v", req)
	var user *model.User
	user, err = receive.userStore.Query(ctx, userstore.ID(req.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierr.InternalServer().WithStack().WithErr(reason.ErrUserNotFound)
		}
		return err
	}
	if *user.Status == model.UserStatusDisable {
		logger.WithContext(ctx, true).Errorf("user has been disabled, user email: %s", user.Email)
		return apierr.InternalServer().WithStack().WithErr(reason.ErrUserNotFound)
	}

	if !receive.verifyPassword(ctx, req.OldPassword, user.Password) {
		return apierr.InternalServer().WithStack().WithErr(reason.ErrInvalidPassword)
	}

	encryptPassword, err := receive.encryptPassword(ctx, req.NewPassword)
	if err != nil {
		return err
	}
	user.Password = encryptPassword

	err = receive.userStore.Save(ctx, user)
	if err != nil {
		return err
	}
	if receive.ldapEnable {
		err = receive.ldap.UpdateUserPassword(ctx, user.Name, req.NewPassword)
		if err != nil {
			return err
		}
	}
	return nil
}

func (receive *UserSVC) UpdateUser(ctx context.Context, req *schema.UserUpdateRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("user update, request: %#v", req)
	var user *model.User
	user, err = receive.userStore.Query(ctx, userstore.ID(req.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierr.InternalServer().WithStack().WithErr(reason.ErrUserNotFound)
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

// UpdateUserRole 更新用户角色
func (receive *UserSVC) UpdateUserRole(ctx context.Context, req *schema.UserUpdateRoleRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("user update role, request: %#v", req)
	var (
		user        *model.User
		oldRoleName string
	)
	user, err = receive.userStore.Query(ctx, userstore.ID(req.ID), userstore.LoadRole())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierr.InternalServer().WithStack().WithErr(reason.ErrUserNotFound)
		}
		return err
	}
	if *user.Status == model.UserStatusDisable {
		logger.WithContext(ctx, true).Errorf("user has been disabled, userName: %s", user.Name)
		return apierr.InternalServer().WithStack().WithErr(reason.ErrUserNotFound)
	}

	if user.RoleName == req.RoleName {
		return apierr.InternalServer().WithStack().WithErr(reason.ErrUserHasSameRole)
	}
	oldRoleName = user.RoleName
	user.RoleName = req.RoleName
	user.Role = nil
	role, err := receive.roleStore.Query(ctx, rbac.RoleName(req.RoleName))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierr.InternalServer().WithStack().WithErr(reason.ErrRoleNotFound)
		}
		return err
	}
	userCache := helpers.GetRoleCacheKey(user.Name)
	if err = receive.cache.Del(ctx, userCache); err != nil {
		return err
	}
	if err = receive.userStore.Save(ctx, user); err != nil {
		return err
	}
	if err = receive.cache.SetString(ctx, userCache, role.Name, &cache.NeverExpires); err != nil {
		return err
	}

	go func() {
		time.Sleep(time.Millisecond * 200)
		if err = receive.cache.Del(ctx, userCache); err != nil {
			logger.WithContext(ctx, true).Error(err)
		}
	}()

	if receive.ldapEnable {
		if oldRoleName != "" {
			err = receive.ldap.RemoveUserFromGroup(ctx, oldRoleName, user.Name)
			if err != nil {
				return err
			}
		}
		exist, err := receive.ldap.SearchGroup(ctx, req.RoleName)
		if err != nil {
			return err
		}
		if !exist {
			err = receive.ldap.CreateGroup(ctx, req.RoleName)
			if err != nil {
				return err
			}
		}
		err = receive.ldap.AddUserToGroup(ctx, req.RoleName, user.Name)
		if err != nil {
			return err
		}
	}
	return
}

// Info 获取用户信息
func (receive *UserSVC) Info(ctx context.Context, id int) (res *schema.UserResponse, err error) {
	logger.WithContext(ctx, true).Debugf("user info, request: %#v", id)
	user, err := receive.userStore.Query(ctx, userstore.ID(id), userstore.LoadRole(), userstore.LoadRolePolicy())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.InternalServer().WithStack().WithErr(reason.ErrUserNotFound)
		}
		return nil, err
	}
	if *user.Status == model.UserStatusDisable {
		logger.WithContext(ctx, true).Errorf("userstore has been disabled, userName: %s", user.Name)
		return nil, apierr.InternalServer().WithStack().WithErr(reason.ErrUserNotFound)
	}
	res = &schema.UserResponse{}
	res.ConvertToUserResponse(user)
	return res, nil
}

func (receive *UserSVC) ListUser(ctx context.Context, req *schema.UserListRequest) (data *schema.UserListResponse, err error) {
	logger.WithContext(ctx, true).Debugf("user list, request: %#v", req)
	total, users, err := receive.userStore.List(ctx, req.Page, req.PageSize)
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
		return "", apierr.InternalServer().WithStack().WithMsg("failed to encrypt password").WithErr(err)
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
