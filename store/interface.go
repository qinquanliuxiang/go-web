package store

import (
	"context"
	"qqlx/model"
	"qqlx/store/rbac"
	"qqlx/store/userstore"
	"time"
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

// PolicyStoreInterface 策略CRUD
type PolicyStoreInterface interface {
	// Query 查询策略
	//
	// @param options 查询选项
	// @return policy 策略
	// @return err 错误
	Query(ctx context.Context, options ...rbac.PolicyQueryOption) (policy *model.Policy, err error)
	Create(ctx context.Context, policy *model.Policy) (err error)
	Save(ctx context.Context, policy *model.Policy) (err error)
	// Delete 删除策略
	//
	// @param policy 策略
	// @param options 删除选项
	// @return err 错误
	Delete(ctx context.Context, policy *model.Policy, options ...rbac.PolicyDeleteOption) (err error)
	List(ctx context.Context, page, pageSize int, options ...rbac.PolicyQueryOption) (total int64, polices []model.Policy, err error)
}

// CacheInterface 缓存
type CacheInterface interface {
	GetSet(ctx context.Context, key string) ([]string, error)

	SetSet(ctx context.Context, key string, value []any, expireTime *time.Duration) error
	// GetString 获取字符串
	//
	// @param key 键
	// @return data 数据
	// @return err 错误
	GetString(ctx context.Context, key string) (data string, err error)
	// SetString 设置字符串
	//
	// @param key 键
	// @param value 值
	// @param expireTime 过期时间
	// @return err 错误
	SetString(ctx context.Context, key, value string, expireTime *time.Duration) (err error)
	// GetInt64 获取整数
	//
	// @param key 键
	// @return data 数据
	// @return err 错误
	GetInt64(ctx context.Context, key string) (data *int64, err error)
	// SetInt64 设置整数
	//
	// @param key 键
	// @param value 值
	// @param expireTime 过期时间
	// @return err 错误
	SetInt64(ctx context.Context, key string, value int64, expireTime *time.Duration) (err error)
	// Incr 自增
	// @param key 键
	// @return int64 自增后的值
	// @return err 错误
	Incr(ctx context.Context, key string) (int64, error)
	// Del 删除
	//
	// @param key 键
	// @return err 错误
	Del(ctx context.Context, key string) (err error)
	// Flush 清空所有的键值对
	//
	// @return err 错误
	Flush(ctx context.Context) (err error)
}

// LdapInterface LDAP接口
type LdapInterface interface {
	LdapUserInterface
	LdapGroupInterface
}

type LdapUserInterface interface {
	// CreateUser 创建用户
	//
	// @param name 用户名
	// @param password 密码
	// @param email 邮箱
	// @return err 错误
	CreateUser(ctx context.Context, name, password, email string) error
	// DeleteUser 删除用户
	//
	// @param username 用户名
	// @return err 错误
	DeleteUser(ctx context.Context, username string) error
	// UpdateUserPassword 更新用户密码
	//
	// @param username 用户名
	// @param password 密码
	// @return err 错误
	UpdateUserPassword(ctx context.Context, username, password string) error
	// SearchUser 查询用户
	//
	// @param username 用户名
	// @return user 用户
	SearchUser(ctx context.Context, username string) (*model.User, error)
	// SearchUserGroups 查询用户组
	//
	// @param username 用户名
	// @return groups 用户组
	// @return err 错误
	SearchUserGroups(_ context.Context, username string) (groups []string, err error)
}

type LdapGroupInterface interface {
	// CreateGroup 创建组
	//
	// @param groupName 组名
	// @return err 错误
	CreateGroup(ctx context.Context, groupName string) error
	// DeleteGroup 删除组
	//
	// @param groupName 组名
	// @return err 错误
	DeleteGroup(ctx context.Context, groupName string) error
	// SearchGroup 查询组
	//
	// @param groupName 组名
	// @return exist 组是否存在
	// @return err 错误
	SearchGroup(ctx context.Context, groupName string) (exist bool, err error)
	// AddUserToGroup 添加用户到组
	//
	// @param groupName 组名
	// @param userName 用户名
	// @return err 错误
	AddUserToGroup(ctx context.Context, groupName, userName string) error
	// RemoveUserFromGroup 从组中移除用户
	//
	// @param groupName 组名
	// @param userName 用户名
	// @return err 错误
	RemoveUserFromGroup(ctx context.Context, groupName, userName string) error
	// SearchGroupMembers 查询组成员
	//
	// @param groupName 组名
	// @return group 组
	// @return err 错误
	SearchGroupMembers(ctx context.Context, groupName string) (group *model.LdapGroup, err error)
}
