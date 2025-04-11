package constant

const (
	FlagConfigPath = "config-path"
	ConfigEnv      = "CONFIG_PATH"
)

const (
	ServerVersion        = "1.0.0"
	DefaultServerBind    = "0.0.0.0:8080"
	DefaultServerName    = "qqlx"
	DefaultJwtExpireTime = "12h"
	DefaultJwtIssuer     = "qqlx"
	DefaultLoglevel      = "info"
	DefaultRedisIncrKey  = "machine_id"
	AuthMidwareKey       = "user"
	LogErrMidwareKey     = "error"
	TraceID              = "traceID"
)

// redis
const (
	// DefaultRedisExpireTime 默认redis过期时间
	DefaultRedisExpireTime = "30s"
	// RoleCacheKeyPrefix RedisKeyPrefix redis 角色缓存 key 前缀
	RoleCacheKeyPrefix = "role"
)
