package interfaces

import (
	"context"
	"qqlx/model"
	"qqlx/store/rbac"
)

// RoleStoreInterface 角色CRUD
type RoleStoreInterface interface {
	// Query 查询角色
	//
	// @param options 查询选项
	// @return role 角色
	// @return err 错误
	Query(ctx context.Context, options ...rbac.RoleQueryOption) (role *model.Role, err error)
	Create(ctx context.Context, role *model.Role) (err error)
	Save(ctx context.Context, role *model.Role) (err error)
	// Delete 删除角色
	//
	// @param role 角色
	// @param options 删除选项
	// @return err 错误
	Delete(ctx context.Context, role *model.Role, options ...rbac.RoleDeleteOption) (err error)
	List(ctx context.Context, page, pageSize int, options ...rbac.RoleQueryOption) (total int64, roles []model.Role, err error)
}

// RolePolicyStoreInterface 角色关联策略
type RolePolicyStoreInterface interface {
	// AppendPolicy 追加策略
	//
	// @param role 角色, 追加策略的角色
	// @param policy 策略, 需要追加的策略
	// @return err 错误
	AppendPolicy(ctx context.Context, role *model.Role, policy []model.Policy) (err error)

	// DeletePolicy 删除策略
	//
	// @param role 角色, 删除策略的角色
	// @param policy 策略, 需要删除的策略
	// @return err 错误
	DeletePolicy(ctx context.Context, role *model.Role, policy []model.Policy) (err error)
}
