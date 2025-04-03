package model

import (
	"gorm.io/plugin/soft_delete"
)

type Policy struct {
	ID        int                   `gorm:"primarykey"`
	CreatedAt int                   `gorm:"autoCreateTime"`
	UpdatedAt int                   `gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:;index"`
	Name      string                `gorm:"comment:名称;size:50;uniqueIndex:idx_policy_name_path_method"`
	Path      string                `gorm:"comment:路径;size:128;uniqueIndex:idx_policy_name_path_method"`
	Method    string                `gorm:"comment:方法;size:10;uniqueIndex:idx_policy_name_path_method"`
	Describe  string                `gorm:"comment:描述;size:1024"`
	Roles     []*Role               `gorm:"many2many:role_policy;"`
}

func (p *Policy) TableName() string {
	return "policys"
}
