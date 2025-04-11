package schema

import (
	"qqlx/model"

	"gorm.io/plugin/soft_delete"
)

type PolicyCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Describe string `json:"describe" validate:"required"`
	Path     string `json:"path" validate:"required"`
	Method   string `json:"method" validate:"required"`
}

type PolicyIDRequest struct {
	ID int `uri:"id" validate:"required,gte=1"`
}

type PolicyUpdateRequest struct {
	ID       int    `uri:"id" validate:"required"`
	Describe string `json:"describe" validate:"required"`
}

type PolicyListRequest struct {
	Page     int `form:"page" validate:"required,gt=0|eq=-1" json:"page"`
	PageSize int `form:"pageSize" validate:"required,gt=0|eq=-1" json:"pageSize"`
}

type PolicyListResponse struct {
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Items    []model.Policy `json:"items"`
}

type PolicyResponse struct {
	ID          int                   `json:"id"`
	CreatedAt   int                   `json:"createdAt"`
	UpdatedAt   int                   `json:"updatedAt"`
	DeletedAt   soft_delete.DeletedAt `json:"deletedAt"`
	Name        string                `json:"name"`
	Path        string                `json:"path"`
	Method      string                `json:"method"`
	Description string                `json:"description"`
	Roles       []model.Role          `json:"roles"`
}
