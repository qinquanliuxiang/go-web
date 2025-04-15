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

// PolicySortByCreatedDesc 按照创建时间倒序
func PolicySortByCreatedDesc() PolicyQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Order("created_at desc")
	}
}

// InPolicy 根据 policy id 列表查询
func InPolicy(ids []int) PolicyQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("id in (?)", ids)
	}
}

// NotInPolicyNames 过滤掉 names 的策略
func NotInPolicyNames(names []string) PolicyQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("name NOT IN ?", names)
	}
}

// PolicyQueryByName 根据 name 进行前缀查询
func PolicyQueryByName(keyword string, value string) PolicyQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		likeVal := value + "%"
		switch keyword {
		case "name":
			query = query.Where("name LIKE ?", likeVal)
		}
		return query
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
		return nil, apierr.InternalServer().Set(apierr.DBErrCode, "failed to query policy", err)
	}
	return policy, nil
}

func (receive *PolicyStore) Create(ctx context.Context, policy *model.Policy) (err error) {
	if err = receive.data.WithContext(ctx).Create(&policy).Error; err != nil {
		return apierr.InternalServer().Set(apierr.DBErrCode, "failed to create policy", err)
	}
	return nil
}

func (receive *PolicyStore) Save(ctx context.Context, policy *model.Policy) (err error) {
	if err = receive.data.WithContext(ctx).Save(&policy).Error; err != nil {
		return apierr.InternalServer().Set(apierr.DBErrCode, "failed to save policy", err)
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
		return apierr.InternalServer().Set(apierr.DBErrCode, "failed to delete policy", err)
	}
	return nil
}

func (receive *PolicyStore) List(ctx context.Context, page int, pageSize int, options ...PolicyQueryOption) (total int64, polices []model.Policy, err error) {
	query := receive.data.WithContext(ctx).Model(&model.Policy{})
	// 添加查询选项
	for _, option := range options {
		query = option(query)
	}

	err = query.Count(&total).Error
	if err != nil {
		return 0, nil, apierr.InternalServer().Set(apierr.DBErrCode, "failed to count policies", err)
	}

	if page == -1 && pageSize == -1 {
		if err = query.Find(&polices).Error; err != nil {
			return 0, nil, apierr.InternalServer().Set(apierr.DBErrCode, "failed to list policies", err)
		}
		return total, polices, nil
	}

	err = query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&polices).Error
	if err != nil {
		return 0, nil, apierr.InternalServer().Set(apierr.DBErrCode, "failed to list policies", err)
	}
	return total, polices, nil
}
