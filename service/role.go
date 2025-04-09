package service

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"qqlx/base/apierr"
	"qqlx/base/conf"
	"qqlx/base/helpers"
	"qqlx/base/logger"
	"qqlx/base/reason"
	"qqlx/model"
	"qqlx/pkg/sonyflake"
	"qqlx/schema"
	"qqlx/store"
	"qqlx/store/rbac"
)

type RoleSVC struct {
	generateID        *sonyflake.GenerateIDStruct
	roleStore         store.RoleStoreInterface
	policyStore       store.PolicyStoreInterface
	appendPolicyStore store.RolePolicyStoreInterface
	casbinStore       store.CasbinInterface
	ldapEnable        bool
	ldap              store.LdapInterface
}

func NewRoleSVC(
	generateID *sonyflake.GenerateIDStruct,
	generalRoleStore store.RoleStoreInterface,
	policyStore store.PolicyStoreInterface,
	appendStore store.RolePolicyStoreInterface,
	casbinStore store.CasbinInterface,
	ldap store.LdapInterface,
) *RoleSVC {
	ldapEnable := conf.GetLdapEnable()
	return &RoleSVC{
		generateID:        generateID,
		roleStore:         generalRoleStore,
		policyStore:       policyStore,
		casbinStore:       casbinStore,
		appendPolicyStore: appendStore,
		ldapEnable:        ldapEnable,
		ldap:              ldap,
	}
}

func (receive *RoleSVC) GetRole(ctx context.Context, req *schema.RoleIDRequest) (role *model.Role, err error) {
	logger.WithContext(ctx, true).Debugf("get role, request: %#v", req)
	role, err = receive.roleStore.Query(ctx, rbac.RoleID(req.ID), rbac.LoadPolices())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.InternalServer().Set(apierr.ServiceErrCode, "role not found", reason.ErrRoleNotFound)
		}
		return nil, err
	}

	return role, nil
}

func (receive *RoleSVC) CreateRole(ctx context.Context, req *schema.RoleCreateRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("create role, request: %#v", req)
	var (
		id    int
		exits bool
	)
	query, err := receive.roleStore.Query(ctx, rbac.RoleName(req.Name))
	if err != nil {
		if !errors.Is(err, reason.ErrRoleNotFound) {
			return err
		}
	}
	if query != nil {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, reason.ErrRoleExists.Error(), reason.ErrRoleExists)
	}

	id, err = receive.generateID.NextID()
	if err != nil {
		return err
	}
	role := &model.Role{
		ID:          id,
		Name:        req.Name,
		Description: req.Describe,
	}
	if receive.ldapEnable {
		exits, err = receive.ldap.SearchGroup(ctx, role.Name)
		if err != nil {
			return err
		}
		if !exits {
			err = receive.ldap.CreateGroup(ctx, role.Name)
			if err != nil {
				return err
			}
		}
	}
	return receive.roleStore.Create(ctx, role)
}

// DeleteRole 删除角色
func (receive *RoleSVC) DeleteRole(ctx context.Context, req *schema.RoleIDRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("delete role, request: %#v", req)
	role, err := receive.roleStore.Query(ctx, rbac.RoleID(req.ID), rbac.LoadUsers())
	if err != nil {
		return err
	}
	if len(role.Users) > 0 {
		var userNames []string
		for _, user := range role.Users {
			userNames = append(userNames, user.Name)
		}
		return apierr.InternalServer().Set(apierr.ServiceErrCode, fmt.Sprintf("role has user: %v", userNames), reason.ErrRoleHasUser)
	}

	if receive.ldapEnable {
		err = receive.ldap.DeleteGroup(ctx, role.Name)
		if err != nil {
			return err
		}
	}

	deleteCasbin := helpers.GetCasbinRole(role.Name, role.Policys)
	err = receive.casbinStore.DeleteRolePolices(ctx, deleteCasbin)
	if err != nil {
		return err
	}
	err = receive.appendPolicyStore.DeletePolicy(ctx, role, role.Policys)
	if err != nil {
		return err
	}
	return receive.roleStore.Delete(ctx, role, rbac.RoleUnscoped())
}

// UpdateRoleDesc 更新角色描述信息
func (receive *RoleSVC) UpdateRoleDesc(ctx context.Context, req *schema.RoleUpdateRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("get role, request: %#v", req)
	role, err := receive.roleStore.Query(ctx, rbac.RoleID(req.ID))
	if err != nil {
		return err
	}

	if role.Description == req.Describe {
		return nil
	}

	role.Description = req.Describe
	return receive.roleStore.Save(ctx, role)
}

// AddByPolicy 增加角色权限
func (receive *RoleSVC) AddByPolicy(ctx context.Context, req *schema.RolePolicyRequest) (err error) {
	logger.WithContext(ctx, false).Debugf("role add policy, request: %#v", req)
	// 去重
	reqPolicesIDs := helpers.Deduplicate(req.PolicyIds)
	// 获取角色
	role, err := receive.roleStore.Query(ctx, rbac.RoleID(req.ID))
	if err != nil {
		return err
	}
	if role.Name == "admin" {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, reason.ErrAdminUserNotAllow.Error(), reason.ErrAdminUserNotAllow)
	}

	_, list, err := receive.policyStore.List(ctx, 1, len(reqPolicesIDs), rbac.InPolicy(reqPolicesIDs))
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, reason.ErrPolicyNotFound.Error(), reason.ErrPolicyNotFound)
	}

	notFound := helpers.FindMissingByID(list, req.PolicyIds)
	if len(notFound) > 0 {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, fmt.Sprintf("policy not found: %v", notFound), reason.ErrPolicyNotFound)
	}

	// role 追加策略
	err = receive.appendPolicyStore.AppendPolicy(ctx, role, list)
	if err != nil {
		return err
	}
	// 更新 casbin 策略
	saveCasbin := helpers.GetCasbinRole(role.Name, list)
	return receive.casbinStore.CreateRolePolices(ctx, saveCasbin)
}

// DeleteByPolicy 删除角色权限
func (receive *RoleSVC) DeleteByPolicy(ctx context.Context, req *schema.RolePolicyRequest) (err error) {
	logger.WithContext(ctx, true).Debugf("create role, request: %#v", req)
	// 去重
	policesID := helpers.Deduplicate(req.PolicyIds)
	if len(policesID) == 0 {
		return nil
	}

	// 获取角色
	role, err := receive.roleStore.Query(ctx, rbac.RoleID(req.ID))
	if err != nil {
		return err
	}
	if role.Name == "admin" {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, reason.ErrAdminUserNotAllow.Error(), reason.ErrAdminUserNotAllow)
	}

	// 获取数据库中的策略
	_, list, err := receive.policyStore.List(ctx, 1, len(policesID), rbac.InPolicy(policesID))
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, reason.ErrPolicyNotFound.Error(), reason.ErrPolicyNotFound)
	}

	notFound := helpers.FindMissingByID(list, req.PolicyIds)
	if len(notFound) > 0 {
		return apierr.InternalServer().Set(apierr.ServiceErrCode, fmt.Sprintf("policy not found: %v", notFound), reason.ErrPolicyNotFound)
	}

	// 删除策略
	if err = receive.appendPolicyStore.DeletePolicy(ctx, role, list); err != nil {
		return err
	}

	// 删除 casbin 策略
	deleteCasbin := helpers.GetCasbinRole(role.Name, list)
	return receive.casbinStore.DeleteRolePolices(ctx, deleteCasbin)
}

func (receive *RoleSVC) ListRole(ctx context.Context, req *schema.RoleListRequest) (data *schema.RoleListResponse, err error) {
	logger.WithContext(ctx, true).Debugf("role list, request: %#v", req)
	options := make([]rbac.RoleQueryOption, 0)
	if req.Keyword != "" {
		options = append(options, rbac.RoleQueryByName(req.Keyword, req.Value))
	}
	options = append(options, rbac.RoleSortByCreatedDesc())

	total, roles, err := receive.roleStore.List(ctx, req.Page, req.PageSize, options...)
	if err != nil {
		return nil, err
	}
	data = &schema.RoleListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Items:    roles,
	}
	return data, nil
}
