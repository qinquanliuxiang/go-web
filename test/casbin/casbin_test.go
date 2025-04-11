package casbin_test

import (
	"qqlx/base/conf"
	"qqlx/base/data"
	"testing"

	"github.com/casbin/casbin/v2"
)

var (
	enforcer *casbin.Enforcer
	err      error
)

func InitCli() {
	err = conf.LoadConfig("../../config.yaml")
	if err != nil {
		panic(err)
	}
	enforcer, err = data.InitCasbin()
	if err != nil {
		panic(err)
	}
}
func TestCasbinCreateRole(t *testing.T) {
	// test.Enforcer.AddPolicy("admin", "*", "*")

	// ok, err := test.Enforcer.AddPolicies([][]string{
	// 	{"admin", "userstore", "create"},
	// 	{"admin", "userstore", "delete"},
	// 	{"admin", "userstore", "update"},
	// 	{"admin", "userstore", "get"},
	// })
	// t.Logf("ok: %v, err: %v", ok, err)
	// policsy, err := test.Enforcer.GetFilteredPolicy(0, "admin")
	// t.Logf("policsy: %v, err: %v", policsy, err)
	ok, err := enforcer.RemovePolicies([][]string{
		{"admin", "userstore", "create"},
		{"admin", "userstore", "delete"},
		{"admin", "userstore", "update"},
		{"admin", "userstore", "get"},
	})
	t.Logf("ok: %v, err: %v", ok, err)
	// test.Enforcer.SavePolicy()
	// ok, err := test.Enforcer.Enforce("qq", "userstore", "create")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if !ok {
	// 	t.Fatal("authentication failed")
	// }
}
