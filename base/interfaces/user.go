package interfaces

import (
	"context"
	"qqlx/model"
	"qqlx/store/userstore"
)

// UserStoreInterface 用户CRUD
type UserStoreInterface interface {
	// Query 查询用户
	//
	// @param options 查询选项
	// @return user 用户
	// @return err 错误
	Query(ctx context.Context, options ...userstore.QueryOption) (user *model.User, err error)
	Create(ctx context.Context, user *model.User) (err error)
	Save(ctx context.Context, user *model.User) (err error)
	// Delete 删除用户
	//
	// @param user 用户
	// @param options 删除选项
	// @return err 错误
	Delete(ctx context.Context, user *model.User, options ...userstore.DeleteOption) (err error)
	// List 查询用户列表
	//
	// @param page 页码
	// @param pageSize 每页数量
	// @param options 查询选项
	// @return total 总数
	// @return users 用户列表
	// @return err 错误
	List(ctx context.Context, page, pageSize int, options ...userstore.QueryOption) (total int64, users []model.User, err error)
}

// UserRoleStoreInterface 用户角色关系
type UserRoleStoreInterface interface {
	// AppendRoles 用户追加角色
	//
	// @param user 用户, 追加角色的用户
	// @param roles 角色
	// @return err 错误
	AppendRoles(ctx context.Context, user *model.User, roles []model.Role) (err error)
	// DeleteRoles 用户删除角色
	//
	// @param user 用户, 删除角色的用户
	// @param roles 角色
	// @return err 错误
	DeleteRoles(ctx context.Context, user *model.User, roles []model.Role) (err error)
}
