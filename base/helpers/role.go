package helpers

import "qqlx/model"

// GetCasbinRole  获取role拥有的权限
//
// policys [][]string{role, path, method}
func GetCasbinRole(roleName string, policys []model.Policy) [][]string {
	save := make([][]string, len(policys))
	for i, policy := range policys {
		save[i] = []string{roleName, policy.Path, policy.Method}
	}
	return save
}
