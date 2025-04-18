package interfaces

import (
	"context"
	"qqlx/model"
)

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
