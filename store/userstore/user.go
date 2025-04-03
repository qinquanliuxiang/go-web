package userstore

import (
	"context"
	"fmt"
	"os/user"
	"qqlx/base/apierr"
	"qqlx/model"

	"gorm.io/gorm"
)

type QueryOption func(query *gorm.DB) *gorm.DB

// ID 根据 user id 查询
func ID(id int) QueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("id = ?", id)
	}
}

// Name 根据 user name 查询
func Name(name string) QueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("name = ?", name)
	}
}

// Email 根据 user email 查询
func Email(email string) QueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("email = ?", email)
	}
}

// LoadRole 用户预加载 Role
func LoadRole() QueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Preload("Role")
	}
}

// LoadRolePolicy 设置 User 预加载 Role.Policy
func LoadRolePolicy() QueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Preload("Role.Policys")
	}
}

type DeleteOption func(query *gorm.DB) *gorm.DB

// Unscoped 永久删除 user
func Unscoped() DeleteOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Unscoped()
	}
}

type Store struct {
	store *gorm.DB
}

func NewUserStore(store *gorm.DB) *Store {
	return &Store{
		store: store,
	}
}

func (receive *Store) Query(ctx context.Context, options ...QueryOption) (user *model.User, err error) {
	sql := receive.store.WithContext(ctx).Model(&user)
	if len(options) > 0 {
		for _, option := range options {
			sql = option(sql)
		}
	}
	err = sql.Take(&user).Error
	if err != nil {
		return nil, apierr.InternalServer().WithMsg("failed to query user").WithErr(err)
	}
	return
}

func (receive *Store) Create(ctx context.Context, user *model.User) (err error) {
	if user == nil {
		return apierr.InternalServer().WithMsg("failed to create user").WithErr(fmt.Errorf("user is nil"))
	}
	err = receive.store.WithContext(ctx).Create(&user).Error
	if err != nil {
		return apierr.InternalServer().WithMsg("failed to create user").WithErr(err)
	}
	return nil
}

func (receive *Store) Delete(ctx context.Context, user *model.User, options ...DeleteOption) (err error) {
	sql := receive.store.WithContext(ctx).Model(&user)
	if len(options) > 0 {
		for _, option := range options {
			sql = option(sql)
		}
	}
	err = sql.Delete(&user).Error
	if err != nil {
		return apierr.InternalServer().WithMsg("failed to delete user").WithErr(err)
	}
	return nil
}

func (receive *Store) Save(ctx context.Context, user *model.User) (err error) {
	if user == nil {
		return apierr.InternalServer().WithMsg("failed to save user").WithErr(fmt.Errorf("user is nil"))
	}
	if err = receive.store.WithContext(ctx).Save(user).Error; err != nil {
		return apierr.InternalServer().WithMsg("failed to save user").WithErr(err)
	}
	return nil
}

func (receive *Store) List(ctx context.Context, page, pageSize int) (total int64, users []model.User, err error) {
	// 计数查询
	query := receive.store.WithContext(ctx).Model(&user.User{})
	if err = query.Count(&total).Error; err != nil {
		return 0, nil, apierr.InternalServer().WithMsg("failed to count users").WithErr(err)

	}

	// 数据查询
	if err = query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		return 0, nil, apierr.InternalServer().WithMsg("failed to list users").WithErr(err)
	}
	return total, users, nil
}
