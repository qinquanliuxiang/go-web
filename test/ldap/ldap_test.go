package ldap_test

import (
	"context"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"log"
	"qqlx/base/conf"
	"qqlx/base/data"
	ldapClient "qqlx/store/ldap"
	"testing"
)

var (
	l         *ldap.Conn
	err       error
	closeFunc func()
	ldapStore *ldapClient.Store
	ctx       = context.Background()
)

func IntLdap() {
	err = conf.LoadConfig("../../config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	l, closeFunc, err = data.InitLdap()
	if err != nil {
		log.Fatal(err)
	}
	ldapStore, err = ldapClient.NewLdapStore(l)
	if err != nil {
		log.Fatal(err)
	}
}

func TestInitLdap(t *testing.T) {
	IntLdap()
}

// 创建组
func TestLdapCreateGroup(t *testing.T) {
	err = ldapStore.CreateGroup(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}
	defer closeFunc()

	err = ldapStore.CreateGroup(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("组创建成功")
}

// 删除组
func TestLdapDeleteGroup(t *testing.T) {
	err = ldapStore.DeleteGroup(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("组删除成功")
}

// 搜索组
func TestLdapSearchGroup(t *testing.T) {
	defer closeFunc()
	groupList, err := ldapStore.SearchGroup(ctx, "ops")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(groupList)
}

// 添加用户到组
func TestLdapAddUserToGroup(t *testing.T) {
	err = ldapStore.AddUserToGroup(ctx, "dev", "huyunfei")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("添加用户到组成功")
}

// 从组中删除用户
func TestLdapRemoveUserFromGroup(t *testing.T) {
	err = ldapStore.RemoveUserFromGroup(ctx, "dev", "huyunfei")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("从组中删除用户成功")
}

// 搜索组中的成员
func TestLdapSearchGroupMembers(t *testing.T) {
	ldapGroup, err := ldapStore.SearchGroupMembers(ctx, "ops")
	if err != nil {
		t.Fatal(err)
	}
	defer closeFunc()
	fmt.Printf("ops组成员: %#v\n", ldapGroup)
}

func TestLdapCreateUser(t *testing.T) {
	err = ldapStore.CreateUser(ctx, "huyunfei", "123456", "huyunfei@qqlx.net")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("用户创建成功")
}

// 删除用户
func TestLdapDeleteUser(t *testing.T) {
	err = ldapStore.DeleteUser(ctx, "huyunfei")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("用户删除成功")
}

// 修改用户
func TestLdapModifyUser(t *testing.T) {
	err = ldapStore.UpdateUserPassword(ctx, "huyunfei", "1234567")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("用户修改成功")
}

// 搜索用户
func TestLdapSearchUser(t *testing.T) {
	user, err := ldapStore.SearchUser(ctx, "huyunfei11")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

func TestUserGroup(t *testing.T) {
	groups, err := ldapStore.SearchUserGroups(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(groups)
}
