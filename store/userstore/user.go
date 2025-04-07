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

// SortByCreatedDesc 按照创建时间倒序
func SortByCreatedDesc() QueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Order("created_at desc")
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

// QueryByNameOrEmail 根据 name 或 email 进行前缀查询
func QueryByNameOrEmail(keyword string, value string) QueryOption {
	return func(query *gorm.DB) *gorm.DB {
		likeVal := value + "%"
		switch keyword {
		case "name":
			query = query.Where("name LIKE ?", likeVal)
		case "email":
			query = query.Where("email LIKE ?", likeVal)
		}
		return query
	}
}

func Status(status int) QueryOption {
	return func(query *gorm.DB) *gorm.DB {
		if status == model.UserStatusAvailable || status == model.UserStatusDisable {
			return query.Where("status = ?", status)
		}
		return query
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

func (receive *Store) Querys(ctx context.Context, options ...QueryOption) (user []model.User, err error) {
	sql := receive.store.WithContext(ctx).Model(&model.User{})
	if len(options) > 0 {
		for _, option := range options {
			sql = option(sql)
		}
	}
	err = sql.Find(&user).Error
	if err != nil {
		return nil, apierr.InternalServer().WithMsg("failed to query users").WithErr(err)
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

func (receive *Store) List(ctx context.Context, page, pageSize int, options ...QueryOption) (int64, []model.User, error) {
	query := receive.store.WithContext(ctx).Model(&user.User{})

	for _, option := range options {
		query = option(query)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, apierr.InternalServer().WithMsg("failed to count users").WithErr(err)
	}

	var users []model.User
	err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error
	if err != nil {
		return 0, nil, apierr.InternalServer().WithMsg("failed to list users").WithErr(err)
	}
	return total, users, nil
}
