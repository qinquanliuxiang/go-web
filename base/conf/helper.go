package conf

import (
	"fmt"
	"qqlx/base/constant"
	"time"

	"github.com/spf13/viper"
)

func GetLogLevel() string {
	logLevel := viper.GetString("server.logLevel")
	if logLevel == "" {
		logLevel = constant.DefaultLoglevel
	}
	return logLevel
}

func GetServerBind() string {
	bind := viper.GetString("server.bind")
	if bind == "" {
		bind = constant.DefaultServerBind
	}
	return bind
}

func GetProjectName() string {
	projectName := viper.GetString("server.projectName")
	if projectName == "" {
		return constant.DefaultServerName
	}
	return projectName
}

func GetSalt() (string, error) {
	salt := viper.GetString("server.salt")
	if salt == "" {
		return "", fmt.Errorf("server.salt is empty")
	}
	return salt, nil
}

func GetCasbinModelPath() (string, error) {
	modlPath := viper.GetString("casbin.modelPath")
	if modlPath == "" {
		return "", fmt.Errorf("casbin.modelPath is empty")
	}
	return modlPath, nil
}

func GetJwtSecret() (string, error) {
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		return "", fmt.Errorf("jwt.secret is empty")
	}
	return secret, nil
}

func GetJwtIssuer() string {
	issuer := viper.GetString("jwt.issuer")
	if issuer == "" {
		return constant.DefaultJwtIssuer
	}
	return issuer
}

func GetJwtExpirationTime() (time.Duration, error) {
	timeOut := viper.GetString("jwt.expirationTime")
	if timeOut == "" {
		timeOut = constant.DefaultJwtExpireTime
	}
	expire, err := time.ParseDuration(timeOut)
	if err != nil {
		return 0, fmt.Errorf("parse jwt.expirationTime failed: %w", err)
	}
	return expire, nil
}

func GetCasbinDsn() (string, error) {
	user := viper.GetString("mysql.username")
	if user == "" {
		return "", fmt.Errorf("mysql.username is empty")
	}
	pas := viper.GetString("mysql.password")
	if pas == "" {
		return "", fmt.Errorf("mysql.password is empty")
	}
	host := viper.GetString("mysql.host")
	if host == "" {
		return "", fmt.Errorf("mysql.host is empty")
	}
	database := viper.GetString("mysql.database")
	if database == "" {
		return "", fmt.Errorf("mysql.database is empty")
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pas, host, database), nil
}

func GetMysqlDsn() (dsn string, err error) {
	user := viper.GetString("mysql.username")
	if user == "" {
		return "", fmt.Errorf("mysql.username is empty")
	}
	pas := viper.GetString("mysql.password")
	if pas == "" {
		return "", fmt.Errorf("mysql.password is empty")
	}
	host := viper.GetString("mysql.host")
	if host == "" {
		return "", fmt.Errorf("mysql.host is empty")
	}
	database := viper.GetString("mysql.database")
	if database == "" {
		return "", fmt.Errorf("mysql.database is empty")
	}
	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&timeout=10000ms",
		user,
		pas,
		host,
		database,
	)
	return dsn, nil
}

func GetRedisPassword() (string, error) {
	password := viper.GetString("redis.password")
	if password == "" {
		return "", fmt.Errorf("redis.password is empty")
	}
	return password, nil
}
func GetRedisMasterName() (string, error) {
	masterName := viper.GetString("redis.sentinel.masterName")
	if masterName == "" {
		return "", fmt.Errorf("redis.sentinel.masterName is empty")
	}
	return masterName, nil
}
func GetRedisSentinelPassword() (string, error) {
	sentPassword := viper.GetString("redis.sentinel.password")
	if sentPassword == "" {
		return "", fmt.Errorf("redis.sentinel.password is empty")
	}
	return sentPassword, nil
}
func GetRedisSentinelHosts() ([]string, error) {
	sentinelHosts := viper.GetStringSlice("redis.sentinel.hosts")
	if len(sentinelHosts) == 0 {
		return nil, fmt.Errorf("redis.sentinel.hosts is empty")
	}
	return sentinelHosts, nil
}

func GetRedisHost() (string, error) {
	host := viper.GetString("redis.host")
	if host == "" {
		return "", fmt.Errorf("redis.host is empty")
	}
	return host, nil
}

func GetRedisDB() int {
	return viper.GetInt("redis.db")
}

func GetRedisMode() string {
	return viper.GetString("redis.mode")
}

func GetRedisExpireTime() (time.Duration, error) {
	expireTime := viper.GetString("redis.expireTime")
	if expireTime == "" {
		expireTime = constant.DefaultRedisExpireTime
	}
	expire, err := time.ParseDuration(expireTime)
	if err != nil {
		return 0, fmt.Errorf("parse redis.expireTime failed: %w", err)
	}
	return expire, nil
}

func GetRedisKeyPrefix() (string, error) {
	prefix := viper.GetString("redis.keyPrefix")
	if prefix == "" {
		return "", fmt.Errorf("redis.keyPrefix is empty")
	}
	return prefix, nil
}

func GetLdapEnable() bool {
	return viper.GetBool("ldap.enable")
}

func GetLdapHost() (string, error) {
	host := viper.GetString("ldap.host")
	if host == "" {
		return "", fmt.Errorf("ldap.host is empty")
	}
	return host, nil
}

func GetLdapRootDN() (string, error) {
	rootDN := viper.GetString("ldap.rootDN")
	if rootDN == "" {
		return "", fmt.Errorf("ldap.rootDN is empty")
	}
	return rootDN, nil
}

func GetLdapRootPassword() (string, error) {
	rootPassword := viper.GetString("ldap.rootPassword")
	if rootPassword == "" {
		return "", fmt.Errorf("ldap.rootPassword is empty")
	}
	return rootPassword, nil
}

func GetLdapUserBase() (string, error) {
	userBase := viper.GetString("ldap.userBase")
	if userBase == "" {
		return "", fmt.Errorf("ldap.userBase is empty")
	}
	return userBase, nil
}

func GetLdapGroupBase() (string, error) {
	groupSearch := viper.GetString("ldap.groupBase")
	if groupSearch == "" {
		return "", fmt.Errorf("ldap.groupBase is empty")
	}
	return groupSearch, nil
}

func GetLdapUserFilter() (string, error) {
	userSearchFilter := viper.GetString("ldap.userSearchFilter")
	if userSearchFilter == "" {
		return "", fmt.Errorf("ldap.userSearchFilter is empty")
	}
	return userSearchFilter, nil
}

func GetLdapGroupFilter() (string, error) {
	groupSearchFilter := viper.GetString("ldap.groupSearchFilter")
	if groupSearchFilter == "" {
		return "", fmt.Errorf("ldap.groupSearchFilter is empty")
	}
	return groupSearchFilter, nil
}
