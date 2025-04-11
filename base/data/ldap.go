package data

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"qqlx/base/conf"
)

func InitLdap() (l *ldap.Conn, close func(), err error) {
	if !conf.GetLdapEnable() {
		return nil, nil, nil
	}

	host, err := conf.GetLdapHost()
	if err != nil {
		return nil, nil, err
	}
	username, err := conf.GetLdapRootDN()
	if err != nil {
		return nil, nil, err
	}
	password, err := conf.GetLdapRootPassword()
	if err != nil {
		return nil, nil, err
	}

	l, err = ldap.DialURL(host)
	if err != nil {
		return nil, nil, fmt.Errorf("connect ldap failed: %w", err)
	}

	// 管理员绑定
	err = l.Bind(username, password)
	if err != nil {
		return nil, nil, fmt.Errorf("bind ldap failed, username: %s, err: %w", username, err)
	}
	return l, func() { _ = l.Close() }, nil
}
