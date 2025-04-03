package rbac

import (
	"context"
	"qqlx/base/apierr"
	"qqlx/model"

	"gorm.io/gorm"
)

type PolicyQueryOption func(query *gorm.DB) *gorm.DB

// LoadRoles 设置预加载 Roles
func LoadRoles() PolicyQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Preload("Roles")
	}
}

// PolicyName 根据 policy name 条件查询
func PolicyName(name string) PolicyQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("name = ?", name)
	}
}

// PolicyID 根据 policy id 条件查询
func PolicyID(id int) PolicyQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("id = ?", id)
	}
}

// InPolicy 根据 policy id 列表查询
func InPolicy(ids []int) PolicyQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("id in (?)", ids)
	}
}

type PolicyDeleteOption func(query *gorm.DB) *gorm.DB

// PolicyUnscoped 永久删除 policy
func PolicyUnscoped() PolicyDeleteOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Unscoped()
	}
}

type PolicyStore struct {
	data *gorm.DB
}

func NewPolicyStore(store *gorm.DB) *PolicyStore {
	return &PolicyStore{
		data: store,
	}
}

func (receive *PolicyStore) Query(ctx context.Context, options ...PolicyQueryOption) (policy *model.Policy, err error) {
	query := receive.data.WithContext(ctx).Model(&policy)
	// 添加查询选项
	for _, option := range options {
		query = option(query)
	}
	if err = query.Take(&policy).Error; err != nil {
		return nil, apierr.InternalServer().WithMsg("failed to query policy").WithErr(err)
	}
	return policy, nil
}

func (receive *PolicyStore) Querys(ctx context.Context, options ...PolicyQueryOption) (polices []model.Policy, err error) {
	query := receive.data.WithContext(ctx).Model(&polices)
	// 添加查询选项
	for _, option := range options {
		query = option(query)
	}
	if err = query.Find(&polices).Error; err != nil {
		return nil, apierr.InternalServer().WithMsg("failed to query policy").WithErr(err)
	}
	return polices, nil
}

func (receive *PolicyStore) Create(ctx context.Context, policy *model.Policy) (err error) {
	if err = receive.data.WithContext(ctx).Create(&policy).Error; err != nil {
		return apierr.InternalServer().WithMsg("failed to create policy").WithErr(err)
	}
	return nil
}

func (receive *PolicyStore) Save(ctx context.Context, policy *model.Policy) (err error) {
	if err = receive.data.WithContext(ctx).Save(&policy).Error; err != nil {
		return apierr.InternalServer().WithMsg("failed to save policy").WithErr(err)
	}
	return nil
}

// Delete 删除记录
// @params options 可选 Unscoped，添加后永久记录，默认软删除
func (receive *PolicyStore) Delete(ctx context.Context, policy *model.Policy, options ...PolicyDeleteOption) (err error) {
	sql := receive.data.WithContext(ctx).Model(&policy)
	if len(options) > 0 {
		for _, option := range options {
			sql = option(sql)
		}
	}
	if err = sql.Delete(&policy).Error; err != nil {
		return apierr.InternalServer().WithMsg("failed to delete policy").WithErr(err)
	}
	return nil
}

func (receive *PolicyStore) List(ctx context.Context, page int, pageSize int) (total int64, polices []model.Policy, err error) {
	query := receive.data.WithContext(ctx).Model(&model.Policy{})
	err = query.Count(&total).Error
	if err != nil {
		return 0, nil, apierr.InternalServer().WithMsg("failed to count policies").WithErr(err)
	}

	err = query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&polices).Error
	if err != nil {
		return 0, nil, apierr.InternalServer().WithMsg("failed to list policies").WithErr(err)
	}
	return total, polices, nil
}
