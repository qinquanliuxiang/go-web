package ldap

import (
	"context"
	"fmt"
	"qqlx/base/apierr"
	"qqlx/base/conf"
	"qqlx/base/logger"
	"qqlx/base/reason"
	"qqlx/model"
	"regexp"

	"github.com/go-ldap/ldap/v3"
)

type Store struct {
	ldap              *ldap.Conn
	rootDN            string
	userBase          string
	groupBase         string
	userSearchFilter  string
	groupSearchFilter string
}

func NewLdapStore(l *ldap.Conn) (*Store, error) {
	rootDN, err := conf.GetLdapRootDN()
	if err != nil {
		return nil, err
	}
	userBase, err := conf.GetLdapUserBase()
	if err != nil {
		return nil, err
	}
	groupBase, err := conf.GetLdapGroupBase()
	if err != nil {
		return nil, err
	}
	userSearchFilter, err := conf.GetLdapUserFilter()
	if err != nil {
		return nil, err
	}
	groupSearchFilter, err := conf.GetLdapGroupFilter()
	if err != nil {
		return nil, err
	}
	return &Store{
		ldap:              l,
		rootDN:            rootDN,
		userBase:          userBase,
		groupBase:         groupBase,
		userSearchFilter:  userSearchFilter,
		groupSearchFilter: groupSearchFilter,
	}, nil
}

// CreateUser 创建用户
func (receive *Store) CreateUser(_ context.Context, name, password, email string) error {
	dn := fmt.Sprintf("uid=%s,%s", name, receive.userBase)
	userReq := ldap.NewAddRequest(dn, nil)
	userReq.Attribute("objectClass", []string{"inetOrgPerson", "organizationalPerson", "person", "top"})
	userReq.Attribute("uid", []string{name})
	userReq.Attribute("cn", []string{name})
	userReq.Attribute("sn", []string{name})
	userReq.Attribute("mail", []string{email})
	userReq.Attribute("displayName", []string{name})
	userReq.Attribute("userPassword", []string{password})
	if err := receive.ldap.Add(userReq); err != nil {
		return apierr.InternalServer().Set(apierr.LdapErrCode, "ldap create user failed", err)
	}
	return nil
}

// DeleteUser 删除用户
func (receive *Store) DeleteUser(_ context.Context, username string) error {
	dn := fmt.Sprintf("uid=%s,%s", username, receive.userBase)
	userReq := ldap.NewDelRequest(dn, nil)
	if err := receive.ldap.Del(userReq); err != nil {
		return apierr.InternalServer().Set(apierr.LdapErrCode, "ldap delete user failed", err)
	}
	return nil
}

// UpdateUserPassword 修改用户
func (receive *Store) UpdateUserPassword(_ context.Context, username, password string) error {
	dn := fmt.Sprintf("uid=%s,%s", username, receive.userBase)
	userReq := ldap.NewModifyRequest(dn, nil)
	if password != "" {
		userReq.Replace("userPassword", []string{password})
	}

	if err := receive.ldap.Modify(userReq); err != nil {
		return apierr.InternalServer().Set(apierr.LdapErrCode, "ldap update user password failed", err)
	}
	return nil
}

// SearchUser 搜索用户
func (receive *Store) SearchUser(_ context.Context, username string) (*model.User, error) {
	dn := fmt.Sprintf("uid=%s,%s", username, receive.userBase)
	searchReq := ldap.NewSearchRequest(
		receive.userBase,
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		dn,
		[]string{"uid", "cn", "sn", "mail", "userPassword"}, nil)
	searchResult, err := receive.ldap.Search(searchReq)
	if err != nil {
		return nil, apierr.InternalServer().Set(apierr.LdapErrCode, "ldap search user failed", err)
	}
	if len(searchResult.Entries) == 0 {
		return nil, apierr.InternalServer().Set(apierr.LdapErrCode, "ldap search user failed", reason.ErrLdapUserNotFound)
	}

	entre := searchResult.Entries[0]
	user := &model.User{}
	user.Name = entre.GetAttributeValue("uid")
	user.Password = entre.GetAttributeValue("userPassword")
	user.Email = entre.GetAttributeValue("mail")
	return user, nil
}

func (receive *Store) SearchUserGroups(_ context.Context, username string) (groups []string, err error) {
	// 过滤条件，查询所有包含该用户 DN 的组
	userDN := fmt.Sprintf("uid=%s,%s", ldap.EscapeFilter(username), receive.userBase)
	filter := fmt.Sprintf("(member=%s)", ldap.EscapeFilter(userDN))
	// 构造搜索请求
	searchReq := ldap.NewSearchRequest(
		receive.groupBase,      // 搜索组的组织单位
		ldap.ScopeWholeSubtree, // 递归搜索
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"cn"}, // 只需要获取组名
		nil,
	)

	searchResult, err := receive.ldap.Search(searchReq)
	if err != nil {
		return nil, apierr.InternalServer().Set(apierr.LdapErrCode, "ldap search user groups failed", err)
	}

	if len(searchResult.Entries) == 0 {
		return nil, apierr.InternalServer().Set(apierr.LdapErrCode, "ldap search user groups failed", reason.ErrLdapGroupNotFound)
	}

	// 解析组名
	groups = make([]string, 0)
	for _, entry := range searchResult.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}

	return groups, nil
}

// CreateGroup 创建组
func (receive *Store) CreateGroup(_ context.Context, groupName string) error {
	groupDN := fmt.Sprintf("cn=%s,%s", groupName, receive.groupBase)
	groupReq := ldap.NewAddRequest(groupDN, nil)
	groupReq.Attribute("objectClass", []string{"groupOfNames", "top"})
	groupReq.Attribute("cn", []string{groupName})
	groupReq.Attribute("member", []string{receive.rootDN})
	if err := receive.ldap.Add(groupReq); err != nil {
		return apierr.InternalServer().Set(apierr.LdapErrCode, "ldap create group failed", err)
	}
	return nil
}

// DeleteGroup 删除组
func (receive *Store) DeleteGroup(_ context.Context, groupName string) error {
	groupDN := fmt.Sprintf("cn=%s,%s", groupName, receive.groupBase)
	groupReq := ldap.NewDelRequest(groupDN, nil)
	if err := receive.ldap.Del(groupReq); err != nil {
		return apierr.InternalServer().Set(apierr.LdapErrCode, "ldap delete group failed", err)
	}
	return nil
}

// SearchGroup 搜索组
func (receive *Store) SearchGroup(ctx context.Context, groupName string) (exist bool, err error) {
	filter := fmt.Sprintf("(cn=%s)", groupName)
	searchReq := ldap.NewSearchRequest(
		receive.groupBase,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"cn", "member"},
		nil,
	)
	searchResult, err := receive.ldap.Search(searchReq)
	if err != nil {
		return false, apierr.InternalServer().Set(apierr.LdapErrCode, "ldap search group failed", err)
	}
	if len(searchResult.Entries) == 0 {
		logger.WithContext(ctx, true).Infof("ldap search group failed, groupName: %s, group not found", groupName)
		return false, nil
	}

	if groupName == searchResult.Entries[0].GetAttributeValue("cn") {
		return true, nil
	}

	return false, nil
}

// AddUserToGroup 添加用户到组
func (receive *Store) AddUserToGroup(_ context.Context, groupName, userName string) error {
	groupDN := fmt.Sprintf("cn=%s,%s", groupName, receive.groupBase)
	userDN := fmt.Sprintf("uid=%s,%s", userName, receive.userBase)
	groupReq := ldap.NewModifyRequest(groupDN, nil)
	groupReq.Add("member", []string{userDN})
	if err := receive.ldap.Modify(groupReq); err != nil {
		return apierr.InternalServer().Set(apierr.LdapErrCode, "ldap add user to group failed", err)
	}
	return nil
}

// RemoveUserFromGroup 从组中删除用户
func (receive *Store) RemoveUserFromGroup(_ context.Context, groupName, userName string) error {
	groupDN := fmt.Sprintf("cn=%s,%s", groupName, receive.groupBase)
	userDN := fmt.Sprintf("uid=%s,%s", userName, receive.userBase)

	// 删除用户
	groupReq := ldap.NewModifyRequest(groupDN, nil)
	groupReq.Delete("member", []string{userDN})
	if err := receive.ldap.Modify(groupReq); err != nil {
		return apierr.InternalServer().Set(apierr.LdapErrCode, "ldap remove user from group failed", err)
	}
	return nil
}

// SearchGroupMembers 搜索组中的成员
// 返回的成员是用户名
func (receive *Store) SearchGroupMembers(_ context.Context, groupName string) (group *model.LdapGroup, err error) {
	filter := fmt.Sprintf("(cn=%s)", groupName)
	searchReq := ldap.NewSearchRequest(
		receive.groupBase,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"cn", "member"},
		nil,
	)
	searchResult, err := receive.ldap.Search(searchReq)
	if err != nil {
		return nil, apierr.InternalServer().Set(apierr.LdapErrCode, "ldap search group failed", err)
	}

	group = &model.LdapGroup{
		Member: make([]string, 0),
	}
	entry := searchResult.Entries[0]
	group.GroupName = entry.GetAttributeValue("cn")
	memberDNs := entry.GetAttributeValues("member")
	// 解析成员 uid
	rex := fmt.Sprintf(`uid=([^,]+),%s`, receive.userBase)
	reg := regexp.MustCompile(rex)
	for _, dn := range memberDNs {
		match := reg.FindStringSubmatch(dn)
		if len(match) > 1 {
			group.Member = append(group.Member, match[1])
		}
	}
	return group, nil
}
