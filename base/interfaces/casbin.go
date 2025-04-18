package interfaces

import (
	"context"
)

// CasbinInterface casbin 权限接口
type CasbinInterface interface {
	// GetRolePolicyByName 获取角色策略
	//
	// @param role 角色名
	// @return polices 策略, polices [][]string{role, path, method}
	// @return err 错误
	GetRolePolicyByName(ctx context.Context, role string) (polices [][]string, err error)
	// CreateRolePolices 创建角色策略
	//
	// @param polices 策略, polices [][]string{role, path, method}
	// @return err 错误
	CreateRolePolices(ctx context.Context, polices [][]string) (err error)
	// DeleteRolePolices 删除角色策略
	//
	// @param polices 策略, polices [][]string{role, path, method}
	// @return err 错误
	DeleteRolePolices(ctx context.Context, polices [][]string) (err error)
	// UpdateRolePolices 更新角色策略
	//
	// @param roleName 角色名
	// @param polices 策略, polices [][]string{role, path, method}
	// @return err 错误
	UpdateRolePolices(ctx context.Context, roleName string, polices [][]string) (err error)
}

// Authorizer 登录时验证用户权限
type Authorizer interface {
	// EnforceWithCtx 验证用户是否具有权限
	//
	// @param sub 用户名
	// @param obj 资源
	// @param act 操作
	EnforceWithCtx(ctx context.Context, sub, obj, act string) (ok bool, err error)
}
