package interfaces

import (
	"context"
	"qqlx/model"
	"qqlx/store/rbac"
)

// PolicyStoreInterface 策略CRUD
type PolicyStoreInterface interface {
	// Query 查询策略
	//
	// @param options 查询选项
	// @return policy 策略
	// @return err 错误
	Query(ctx context.Context, options ...rbac.PolicyQueryOption) (policy *model.Policy, err error)
	Create(ctx context.Context, policy *model.Policy) (err error)
	Save(ctx context.Context, policy *model.Policy) (err error)
	// Delete 删除策略
	//
	// @param policy 策略
	// @param options 删除选项
	// @return err 错误
	Delete(ctx context.Context, policy *model.Policy, options ...rbac.PolicyDeleteOption) (err error)
	List(ctx context.Context, page, pageSize int, options ...rbac.PolicyQueryOption) (total int64, polices []model.Policy, err error)
}
