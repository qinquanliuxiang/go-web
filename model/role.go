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
	Policys     []Policy              `gorm:"many2many:role_policy;"`
	Users       []User                `gorm:"foreignKey:RoleName;references:Name" json:"users,omitempty"`
}

func (receiver *Role) TableName() string {
	return "roles"
}
