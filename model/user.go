package model

import (
	"gorm.io/plugin/soft_delete"
)

var (
	UserStatusAvailable = 1
	UserStatusDisable   = 2
)

type User struct {
	ID        int                   `gorm:"primarykey"`
	CreatedAt int                   `gorm:"autoCreateTime"`
	UpdatedAt int                   `gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:;index"`
	Name      string                `gorm:"comment:用户名称;uniqueIndex;size:50"`
	NickName  string                `gorm:"comment:用户昵称;size:50"`
	Email     string                `gorm:"comment:邮箱;uniqueIndex;size:100"`
	RoleName  string                `gorm:"comment:用户角色名称;size:50"`
	Password  string                `gorm:"comment:用户密码;size:255"`
	Avatar    string                `gorm:"comment:用户头像;size:1024"`
	Mobile    string                `gorm:"comment:用户手机号;size:20"`
	Status    *int                  `gorm:"comment:用户状态,1可用,2删除;size:1;default:1"`
	Role      *Role                 `gorm:"foreignKey:RoleName;references:Name"`
}

func (u *User) TableName() string {
	return "users"
}
