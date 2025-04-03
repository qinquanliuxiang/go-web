package schema

import (
	"gorm.io/plugin/soft_delete"
	"qqlx/model"
)

type UserIDRequest struct {
	ID int `uri:"id" validate:"required,gte=1"`
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
	RoleName  string                `json:"roleName"`
	Status    int                   `json:"status"`
	Role      *model.Role           `json:"role,omitempty"`
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
	receive.RoleName = in.RoleName
	receive.Status = *in.Status
	receive.Role = in.Role
}

type UserListRequest struct {
	Page     int `form:"page" validate:"required,gt=0"`
	PageSize int `form:"pageSize" validate:"required,gt=0"`
}

type UserListResponse struct {
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Items    []UserResponse `json:"items"`
}

type UserUpdateRoleRequest struct {
	ID       int    `uri:"id" validate:"required"`
	RoleName string `uri:"roleName" validate:"required"`
}
