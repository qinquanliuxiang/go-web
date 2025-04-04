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
	role, err = receive.roleStore.Query(ctx, rbac.RoleID(req.ID), rbac.LoadUsers(), rbac.LoadPolices())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.InternalServer().WithMsg("role not found").WithErr(err)
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
		return apierr.InternalServer().WithMsg(fmt.Sprintf("create role name %s failed", req.Name)).WithErr(reason.ErrRoleExists).WithStack()
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
	role, err := receive.roleStore.Query(ctx, rbac.RoleID(req.ID), rbac.LoadPolices(), rbac.LoadUsers())
	if err != nil {
		return err
	}
	if len(role.Users) > 0 {
		var userNames []string
		for _, user := range role.Users {
			userNames = append(userNames, user.Name)
		}
		return apierr.InternalServer().WithMsg("failed to delete role").WithErr(fmt.Errorf("role id %s has users %s", role.Name, userNames))
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
		return apierr.InternalServer().WithErr(reason.ErrAdminRole).WithStack()
	}

	reqPolices, err := receive.policyStore.Querys(ctx, rbac.InPolicy(reqPolicesIDs))
	if err != nil {
		return err
	}
	if len(reqPolices) == 0 {
		return apierr.InternalServer().WithMsg(fmt.Sprintf("policy %v not exists", reqPolicesIDs)).WithErr(reason.ErrPolicyNotFound)
	}
	if len(reqPolices) != len(reqPolicesIDs) {
		// 获取不存在的策略
		dbPoliceIDs := make([]int, len(reqPolices))
		for i, p := range reqPolices {
			dbPoliceIDs[i] = p.ID
		}
		notExistsIds := helpers.GetUnique(dbPoliceIDs, reqPolicesIDs)
		return apierr.InternalServer().WithMsg(fmt.Sprintf("policy %v not exists", notExistsIds)).WithErr(fmt.Errorf("policy not exists"))
	}
	// role 追加策略
	err = receive.appendPolicyStore.AppendPolicy(ctx, role, reqPolices)
	if err != nil {
		return err
	}
	// 更新 casbin 策略
	save := helpers.GetCasbinRole(role.Name, reqPolices)
	return receive.casbinStore.CreateRolePolices(ctx, save)
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
		return apierr.InternalServer().WithErr(reason.ErrAdminRole).WithStack()
	}

	// 获取数据库中的策略
	reqPolices, err := receive.policyStore.Querys(ctx, rbac.InPolicy(policesID))
	if err != nil {
		return err
	}
	if len(reqPolices) == 0 {
		return apierr.InternalServer().WithMsg(fmt.Sprintf("policy %v not exists", policesID)).WithErr(reason.ErrPolicyNotFound).WithStack()
	}
	// 如果获取到的策略和请求的策略不一致, 那么返回不存在的策略
	if len(reqPolices) != len(policesID) {
		reqPoliceIds := make([]int, len(reqPolices))
		for i := range reqPolices {
			reqPoliceIds[i] = reqPolices[i].ID
		}
		notExistsIds := helpers.GetUnique(reqPoliceIds, policesID)
		return apierr.InternalServer().WithMsg(fmt.Sprintf("policy %v not exists", notExistsIds)).WithErr(fmt.Errorf("policy not exists"))
	}

	// 删除策略
	if err = receive.appendPolicyStore.DeletePolicy(ctx, role, reqPolices); err != nil {
		return err
	}

	// 删除 casbin 策略
	deleteRole := helpers.GetCasbinRole(role.Name, reqPolices)
	return receive.casbinStore.DeleteRolePolices(ctx, deleteRole)
}

func (receive *RoleSVC) ListRole(ctx context.Context, req *schema.RoleListRequest) (data *schema.RoleListResponse, err error) {
	logger.WithContext(ctx, true).Debugf("role list, request: %#v", req)
	total, roles, err := receive.roleStore.List(ctx, req.Page, req.PageSize)
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
