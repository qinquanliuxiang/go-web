package rbac

import (
	"context"
	"errors"
	"qqlx/base/apierr"
	"qqlx/base/reason"
	"qqlx/model"

	"gorm.io/gorm"
)

type RoleQueryOption func(query *gorm.DB) *gorm.DB

// RoleName 根据 role name 查询 role
func RoleName(name string) RoleQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("name = ?", name)
	}
}

func RoleNames(names []string) RoleQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("name in (?)", names)
	}
}

// RoleID  根据 role id 查询 role
func RoleID(id int) RoleQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("id = ?", id)
	}
}

// RoleSortByCreatedDesc 按照创建时间倒序
func RoleSortByCreatedDesc() RoleQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Order("created_at desc")
	}
}

// LoadPolices role 设置预加载 Policy's
func LoadPolices() RoleQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Preload("Policys")
	}
}

// LoadPolicies role 设置预加载 Policy's
func LoadPolicies() RoleQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Preload("Policys")
	}
}

// LoadUsers role 设置预加载 Users
func LoadUsers() RoleQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Preload("Users", "status = 1")
	}
}

// RoleQueryByName 根据 name 进行前缀查询
func RoleQueryByName(keyword string, value string) RoleQueryOption {
	return func(query *gorm.DB) *gorm.DB {
		likeVal := value + "%"
		switch keyword {
		case "name":
			query = query.Where("name LIKE ?", likeVal)
		}
		return query
	}
}

type RoleDeleteOption func(query *gorm.DB) *gorm.DB

// RoleUnscoped 永久删除 role
func RoleUnscoped() RoleDeleteOption {
	return func(query *gorm.DB) *gorm.DB {
		return query.Unscoped()
	}
}

type RoleStore struct {
	store *gorm.DB
}

func NewRoleStore(store *gorm.DB) *RoleStore {
	return &RoleStore{
		store: store,
	}
}

func (receive *RoleStore) Create(ctx context.Context, role *model.Role) (err error) {
	if role == nil {
		return apierr.InternalServer().Set(apierr.DBErrCode, "failed create role", reason.ErrRoleIsEmpty)
	}
	err = receive.store.WithContext(ctx).Create(&role).Error
	if err != nil {
		return apierr.InternalServer().Set(apierr.DBErrCode, "failed create role", err)
	}
	return nil
}

func (receive *RoleStore) Save(ctx context.Context, role *model.Role) (err error) {
	if role == nil {
		return apierr.InternalServer().Set(apierr.DBErrCode, "failed save role", reason.ErrRoleIsEmpty)
	}
	err = receive.store.WithContext(ctx).Save(&role).Error
	if err != nil {
		return apierr.InternalServer().Set(apierr.DBErrCode, "failed save role", err)
	}
	return nil
}

func (receive *RoleStore) Delete(ctx context.Context, role *model.Role, options ...RoleDeleteOption) (err error) {
	sql := receive.store.WithContext(ctx).Model(&role)
	if len(options) > 0 {
		for _, option := range options {
			sql = option(sql)
		}
	}
	err = sql.Delete(&role).Error
	if err != nil {
		return apierr.InternalServer().Set(apierr.DBErrCode, "failed delete role", err)
	}
	return nil
}

func (receive *RoleStore) List(ctx context.Context, page int, pageSize int, options ...RoleQueryOption) (total int64, roles []model.Role, err error) {
	query := receive.store.WithContext(ctx).Model(&model.Role{})

	// 添加查询选项
	for _, option := range options {
		query = option(query)
	}

	// 计数查询
	err = query.Count(&total).Error
	if err != nil {
		return 0, nil, apierr.InternalServer().Set(apierr.DBErrCode, "failed count role", err)
	}

	if page == -1 && pageSize == -1 {
		err = query.Find(&roles).Error
		if err != nil {
			return 0, nil, apierr.InternalServer().Set(apierr.DBErrCode, "failed list role", err)
		}
		return total, roles, nil
	}

	// 数据查询
	err = query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&roles).Error
	if err != nil {
		return 0, nil, apierr.InternalServer().Set(apierr.DBErrCode, "failed list role", err)
	}

	return total, roles, nil
}

func (receive *RoleStore) Query(ctx context.Context, options ...RoleQueryOption) (role *model.Role, err error) {
	query := receive.store.WithContext(ctx).Model(&role)
	// 添加查询选项
	for _, option := range options {
		query = option(query)
	}
	if err = query.Take(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.InternalServer().Set(apierr.DBErrCode, "failed query role", reason.ErrRoleNotFound)
		}
		return nil, apierr.InternalServer().Set(apierr.DBErrCode, "failed query role", err)
	}
	return role, nil
}

type RoleAssociationStore struct {
	store *gorm.DB
}

func NewRoleAssociationStore(store *gorm.DB) *RoleAssociationStore {
	return &RoleAssociationStore{
		store: store,
	}
}

func (r *RoleAssociationStore) AppendPolicy(ctx context.Context, role *model.Role, appendPolicy []model.Policy) (err error) {
	err = r.store.WithContext(ctx).Model(&role).Association("Policys").Append(&appendPolicy)
	if err != nil {
		return apierr.InternalServer().Set(apierr.DBErrCode, "failed append policy", err)
	}
	return nil
}

func (r *RoleAssociationStore) DeletePolicy(ctx context.Context, role *model.Role, policy []model.Policy) (err error) {
	err = r.store.WithContext(ctx).Model(&role).Association("Policys").Delete(&policy)
	if err != nil {
		return apierr.InternalServer().Set(apierr.DBErrCode, "failed delete policy", err)
	}
	return nil
}
