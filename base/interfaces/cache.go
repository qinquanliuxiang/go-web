package interfaces

import (
	"context"
	"time"
)

// CacheInterface 缓存
type CacheInterface interface {
	GetSet(ctx context.Context, key string) ([]string, error)

	SetSet(ctx context.Context, key string, value []any, expireTime *time.Duration) error
	// GetString 获取字符串
	//
	// @param key 键
	// @return data 数据
	// @return err 错误
	GetString(ctx context.Context, key string) (data string, err error)
	// SetString 设置字符串
	//
	// @param key 键
	// @param value 值
	// @param expireTime 过期时间
	// @return err 错误
	SetString(ctx context.Context, key, value string, expireTime *time.Duration) (err error)
	// GetInt64 获取整数
	//
	// @param key 键
	// @return data 数据
	// @return err 错误
	GetInt64(ctx context.Context, key string) (data *int64, err error)
	// SetInt64 设置整数
	//
	// @param key 键
	// @param value 值
	// @param expireTime 过期时间
	// @return err 错误
	SetInt64(ctx context.Context, key string, value int64, expireTime *time.Duration) (err error)
	// Incr 自增
	// @param key 键
	// @return int64 自增后的值
	// @return err 错误
	Incr(ctx context.Context, key string) (int64, error)
	// Del 删除
	//
	// @param key 键
	// @return err 错误
	Del(ctx context.Context, key string) (err error)
	// Flush 清空所有的键值对
	//
	// @return err 错误
	Flush(ctx context.Context) (err error)
}
