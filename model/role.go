package model

import (
	"gorm.io/plugin/soft_delete"
)

type Role struct {
	ID          int                   `gorm:"primarykey" json:"id"`
	CreatedAt   int                   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   int                   `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   soft_delete.DeletedAt `gorm:"softDelete:;index" json:"deletedAt"`
	Name        string                `gorm:"comment:角色名称;uniqueIndex;size:50" json:"name"`
	Description string                `gorm:"comment:角色描述;size:1024" json:"description"`
	Policys     []Policy              `gorm:"many2many:role_policy;" json:"policys,omitempty"`
	Users       []User                `gorm:"many2many:user_role;" json:"users,omitempty"`
}

func (receiver *Role) TableName() string {
	return "roles"
}

func (receiver Role) GetName() string {
	return receiver.Name
}
