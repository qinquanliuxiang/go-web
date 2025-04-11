package schema

import (
	"qqlx/model"
)

type RoleIDRequest struct {
	ID int `uri:"id" validate:"required"`
}

type RoleCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Describe string `json:"desc" validate:"required"`
}

type RoleUpdateRequest struct {
	ID       int    `uri:"id" validate:"required"`
	Describe string `json:"describe" validate:"required"`
}
type RolePolicyRequest struct {
	ID        int   `uri:"id" validate:"required"`
	PolicyIds []int `json:"policyIds" validate:"required"`
}

type RoleListRequest struct {
	Page     int    `form:"page" validate:"required,gt=0|eq=-1"`
	PageSize int    `form:"pageSize" validate:"required,gt=0|eq=-1"`
	Keyword  string `form:"keyword" validate:"omitempty,oneof=name"` // 支持 name 或 email 前缀模糊搜索
	Value    string `form:"value" validate:"required_with=Keyword"`
}

type RoleListResponse struct {
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
	Items    []model.Role `json:"items"`
}
