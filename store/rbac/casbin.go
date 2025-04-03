package rbac

import (
	"context"
	"fmt"
	"qqlx/base/apierr"

	"github.com/casbin/casbin/v2"
)

type CasbinStore struct {
	enforcer *casbin.Enforcer
}

func NewCasbinStore(enforcer *casbin.Enforcer) *CasbinStore {
	return &CasbinStore{
		enforcer: enforcer,
	}
}

// GetRolePolicyByName  GetRolePolicyByName 根据role获取权限
func (receive *CasbinStore) GetRolePolicyByName(_ context.Context, role string) (policys [][]string, err error) {
	policys, err = receive.enforcer.GetFilteredPolicy(0, role)
	if err != nil {
		return nil, apierr.InternalServer().WithMsg(fmt.Sprintf("failed to get casbin role: %s", role)).WithErr(err)
	}
	return policys, nil
}

// CreateRolePolices CreateRolePolicy 创建role拥有的权限
//
// polices [][]string{role, path, method}
func (receive *CasbinStore) CreateRolePolices(_ context.Context, polices [][]string) (err error) {
	_, err = receive.enforcer.AddPolicies(polices)
	if err != nil {
		return apierr.InternalServer().WithMsg("failed to create casbin policy").WithErr(err)
	}

	return nil
}

// DeleteRolePolices 删除role拥有的权限
//
// polices [][]string{role, path, method}
func (receive *CasbinStore) DeleteRolePolices(_ context.Context, polices [][]string) (err error) {
	_, err = receive.enforcer.RemovePolicies(polices)
	if err != nil {
		return apierr.InternalServer().WithMsg("failed to delete casbin policy").WithErr(err)
	}
	return nil
}

// UpdateRolePolices  更新role拥有的权限
//
// polices [][]string{role, path, method}
func (receive *CasbinStore) UpdateRolePolices(ctx context.Context, roleName string, polices [][]string) (err error) {
	oldPolicys, err := receive.GetRolePolicyByName(ctx, roleName)
	if err != nil {
		return err
	}

	_, err = receive.enforcer.UpdatePolicies(oldPolicys, polices)
	if err != nil {
		return apierr.InternalServer().WithMsg("failed to update casbin policy").WithErr(err)
	}
	return nil
}

type Authentication struct {
	enforcer *casbin.Enforcer
}

func NewAuthentication(enforcer *casbin.Enforcer) *Authentication {
	return &Authentication{
		enforcer: enforcer,
	}
}

func (a *Authentication) EnforceWithCtx(_ context.Context, sub, obj, act string) (ok bool, err error) {
	ok, err = a.enforcer.Enforce(sub, obj, act)
	if err != nil {
		return false, apierr.Forbidden().WithErr(err).WithMsg("failed to enforce")
	}
	return ok, nil
}
