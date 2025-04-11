package model

import (
	"gorm.io/plugin/soft_delete"
)

type Policy struct {
	ID        int                   `gorm:"primarykey" json:"id"`
	CreatedAt int                   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt int                   `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:;index" json:"deletedAt"`
	Name      string                `gorm:"comment:名称;size:50;uniqueIndex:idx_policy_name_path_method" json:"name"`
	Path      string                `gorm:"comment:路径;size:128;uniqueIndex:idx_policy_name_path_method" json:"path"`
	Method    string                `gorm:"comment:方法;size:10;uniqueIndex:idx_policy_name_path_method" json:"method"`
	Describe  string                `gorm:"comment:描述;size:1024" json:"describe"`
	Roles     []Role                `gorm:"many2many:role_policy;" json:"roles,omitempty"`
}

func (receiver *Policy) TableName() string {
	return "policys"
}

func (receiver Policy) GetID() int {
	return receiver.ID
}
