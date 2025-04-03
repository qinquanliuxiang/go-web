package service

import (
	"context"
	"errors"
	"fmt"
	"qqlx/base/apierr"
	"qqlx/base/logger"
	"qqlx/base/reason"
	"qqlx/model"
	"qqlx/pkg/sonyflake"
	"qqlx/schema"
	"qqlx/store"
	"qqlx/store/rbac"

	"gorm.io/gorm"
)

type PolicySVC struct {
	generateID  *sonyflake.GenerateIDStruct
	policyStore store.PolicyStoreInterface
}

func NewPolicySVC(generateID *sonyflake.GenerateIDStruct, policyStore store.PolicyStoreInterface) *PolicySVC {
	return &PolicySVC{
		generateID:  generateID,
		policyStore: policyStore,
	}
}

func (receive *PolicySVC) GetPolicy(ctx context.Context, req *schema.PolicyIDRequest) (res *model.Policy, err error) {
	logger.WithContext(ctx, true).Debugf("get policy, request: %#v", req)
	res, err = receive.policyStore.Query(ctx, rbac.PolicyID(req.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.InternalServer().WithMsg("policy not found").WithErr(err)
		}
		return nil, err
	}
	return res, nil
}

func (receive *PolicySVC) CreatePolicy(ctx context.Context, req *schema.PolicyCreateRequest) (err error) {
	logger.WithContext(ctx, false).Debugf("create policy, request: %#v", req)
	id, err := receive.generateID.NextID()
	if err != nil {
		return err
	}
	return receive.policyStore.Create(ctx, &model.Policy{
		ID:       id,
		Name:     req.Name,
		Path:     req.Path,
		Method:   req.Method,
		Describe: req.Describe,
	})
}

// DeletePolicy 删除策略
func (receive *PolicySVC) DeletePolicy(ctx context.Context, req *schema.PolicyIDRequest) (err error) {
	logger.WithContext(ctx, false).Debugf("get policy, request: %#v", req)
	policy, err := receive.policyStore.Query(ctx, rbac.PolicyID(req.ID), rbac.LoadRoles())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierr.InternalServer().WithErr(reason.ErrPolicyNotFound).WithStack().WithMsg(fmt.Sprintf("policy id: %s does not exist ", req.ID))
		}
		return err
	}

	if len(policy.Roles) > 0 {
		var roleNames []string
		for _, role := range policy.Roles {
			roleNames = append(roleNames, role.Name)
		}
		return apierr.InternalServer().WithMsg("failed to delete policy").WithErr(fmt.Errorf("policy %s belongs roles %s", policy.Name, roleNames))
	}

	return receive.policyStore.Delete(ctx, policy, rbac.PolicyUnscoped())
}

// UpdatePolicy 更新策略描述信息
func (receive *PolicySVC) UpdatePolicy(ctx context.Context, req *schema.PolicyUpdateRequest) (err error) {
	logger.WithContext(ctx, false).Debugf("get policy, request: %#v", req)
	policy, err := receive.policyStore.Query(ctx, rbac.PolicyID(req.ID))
	if err != nil {
		return err
	}
	if policy.Describe == req.Describe {
		return nil
	}

	policy.Describe = req.Describe
	return receive.policyStore.Save(ctx, policy)
}

func (receive *PolicySVC) List(ctx context.Context, req *schema.PolicyListRequest) (res *schema.PolicyListResponse, err error) {
	logger.WithContext(ctx, false).Debugf("policy list, request: %#v", req)
	total, polices, err := receive.policyStore.List(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	res = &schema.PolicyListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Items:    polices,
	}
	return res, nil
}
