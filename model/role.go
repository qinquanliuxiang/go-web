package model

import (
	"gorm.io/plugin/soft_delete"
)

type Role struct {
	ID          int                   `gorm:"primarykey"`
	CreatedAt   int                   `gorm:"autoCreateTime"`
	UpdatedAt   int                   `gorm:"autoUpdateTime"`
	DeletedAt   soft_delete.DeletedAt `gorm:"softDelete:;index"`
	Name        string                `gorm:"comment:角色名称;uniqueIndex;size:50"`
	Description string                `gorm:"comment:角色描述;size:1024"`
	Policys     []Policy              `gorm:"many2many:role_policy;" json:"policys,omitempty"`
	Users       []User                `gorm:"many2many:user_role;" json:"users,omitempty"`
}

func (receiver *Role) TableName() string {
	return "roles"
}

func (receiver Role) GetName() string {
	return receiver.Name
}
