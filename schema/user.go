package schema

import (
	"qqlx/model"

	"gorm.io/plugin/soft_delete"
)

type UserQueryRequest struct {
	ID    int      `uri:"id" validate:"required,gte=1"`
	Query []string `form:"query"`
}

type UserEnableRequest struct {
	ID       int    `uri:"id" validate:"required,gte=1"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserListRequest struct {
	Page     int    `form:"page" validate:"required,gt=0"`
	PageSize int    `form:"pageSize" validate:"required,gt=0"`
	Status   int    `form:"status" validate:"required,oneof=-1 1 2"`       // -1:全部 1:启用 2:禁用
	Keyword  string `form:"keyword" validate:"omitempty,oneof=name email"` // 支持 name 或 email 前缀模糊搜索
	Value    string `form:"value" validate:"required_with=Keyword"`        // keyword存在的时候Value一定要存在
}

type UserNameRequest struct {
	Name string `uri:"name" validate:"required,gte=1"`
}

type UserRegistryRequest struct {
	Name     string `json:"name" validate:"required"`
	NickName string `json:"nickName"`
	Password string `json:"password" validate:"required,min=8"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email" validate:"email"`
	Mobile   string `json:"mobile"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserUpdatePasswordRequest struct {
	ID          int
	OldPassword string `json:"oldPassword" validate:"required,min=8"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

type UserUpdateRequest struct {
	ID       int
	NickName string `json:"nickName"`
	Avatar   string `json:"avatar"`
	Mobile   string `json:"mobile"`
}

type UserResponse struct {
	ID        int                   `json:"id"`
	CreatedAt int                   `json:"createdAt"`
	UpdatedAt int                   `json:"updatedAt"`
	DeletedAt soft_delete.DeletedAt `json:"deletedAt"`
	Name      string                `json:"name"`
	NickName  string                `json:"nickName"`
	Avatar    string                `json:"avatar"`
	Email     string                `json:"email"`
	Mobile    string                `json:"mobile"`
	Status    int                   `json:"status"`
	RoleName  []string              `json:"roleName,omitempty"`
	Roles     []model.Role          `json:"roles,omitempty"`
}

func (receive *UserResponse) ConvertToUserResponse(in *model.User) {
	receive.ID = in.ID
	receive.CreatedAt = in.CreatedAt
	receive.UpdatedAt = in.UpdatedAt
	receive.DeletedAt = in.DeletedAt
	receive.Name = in.Name
	receive.NickName = in.NickName
	receive.Avatar = in.Avatar
	receive.Email = in.Email
	receive.Mobile = in.Mobile
	receive.Status = *in.Status
	receive.Roles = in.Roles
}

type UserListResponse struct {
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Items    []UserResponse `json:"items"`
}

type UserUpdateRoleRequest struct {
	ID        int      `uri:"id" validate:"required"`
	RoleNames []string `json:"roleNames" validate:"required"`
}
