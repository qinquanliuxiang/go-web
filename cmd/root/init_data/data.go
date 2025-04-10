package init_data

import (
	"qqlx/schema"
)

var polices = []schema.PolicyCreateRequest{
	{
		Name:     "admin",
		Path:     "*",
		Method:   "*",
		Describe: "超级管理员",
	},
	{
		Name:     "view",
		Path:     "*",
		Method:   "GET",
		Describe: "查看",
	},
	{
		Name:     "listUser",
		Path:     "/api/v1/users",
		Method:   "GET",
		Describe: "获取用户列表",
	},
	{
		Name:     "userInfo",
		Path:     "/api/v1/users/:id",
		Method:   "GET",
		Describe: "根据id获取用户信息",
	},
	{
		Name:     "disableUser",
		Path:     "/api/v1/users/:id",
		Method:   "DELETE",
		Describe: "禁用用户",
	},
	{
		Name:     "enableUser",
		Path:     "/api/v1/users/enable/:id",
		Method:   "PUT",
		Describe: "启用用户",
	},
	{
		Name:     "AddUserRole",
		Path:     "/api/v1/users/:id/roles",
		Method:   "PUT",
		Describe: "增加用户角色",
	},
	{
		Name:     "DeleteUserRole",
		Path:     "/api/v1/users/:id/roles",
		Method:   "POST",
		Describe: "删除用户角色",
	},
	{
		Name:     "listRoles",
		Path:     "/api/v1/roles",
		Method:   "GET",
		Describe: "获取角色列表",
	},
	{
		Name:     "createRole",
		Path:     "/api/v1/roles",
		Method:   "POST",
		Describe: "创建角色",
	},
	{
		Name:     "updateRoleInfo",
		Path:     "/api/v1/roles/:id",
		Method:   "PUT",
		Describe: "更新角色信息",
	},
	{
		Name:     "getRoleInfo",
		Path:     "/api/v1/roles/:id",
		Method:   "GET",
		Describe: "获取角色信息",
	},
	{
		Name:     "deleteRole",
		Path:     "/api/v1/roles/:id",
		Method:   "DELETE",
		Describe: "删除角色",
	},
	{
		Name:     "deleteRole",
		Path:     "/api/v1/role/deleteRole",
		Method:   "POST",
		Describe: "删除角色",
	},
	{
		Name:     "updateRolePolices",
		Path:     "/api/v1/roles/:id/polices",
		Method:   "POST",
		Describe: "增加角色的权限",
	},
	{
		Name:     "deleteRolePolices",
		Path:     "/api/v1/roles/:id/polices",
		Method:   "DELETE",
		Describe: "删除角色权限",
	},
	{
		Name:     "policesList",
		Path:     "/api/v1/polices",
		Method:   "GET",
		Describe: "获取策略列表",
	},
	{
		Name:     "createPolicy",
		Path:     "/api/v1/polices",
		Method:   "POST",
		Describe: "创建策略",
	},
	{
		Name:     "getPolicyInfo",
		Path:     "/api/v1/polices/:id",
		Method:   "GET",
		Describe: "获取策略信息",
	},
	{
		Name:     "updatePolicy",
		Path:     "/api/v1/polices/:id",
		Method:   "PUT",
		Describe: "修改策略信息",
	},
	{
		Name:     "deletePolicy",
		Path:     "/api/v1/polices/:id",
		Method:   "POST",
		Describe: "删除策略",
	},
}
