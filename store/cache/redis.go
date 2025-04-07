package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"qqlx/base/apierr"
	"qqlx/base/conf"
	"time"
)

var (
	NeverExpires time.Duration = 0
)

// Store redis 客户端
type Store struct {
	client     *redis.Client
	expireTime time.Duration
	keyPrefix  string
}

func NewStore(client *redis.Client) (*Store, func(), error) {
	expireTime, err := conf.GetRedisExpireTime()
	if err != nil {
		return nil, nil, err
	}
	closeup := func() {
		_ = client.Close()
	}
	prefix, err := conf.GetRedisKeyPrefix()
	if err != nil {
		return nil, nil, err
	}
	return &Store{
		client:     client,
		expireTime: expireTime,
		keyPrefix:  prefix,
	}, closeup, nil
}

func (c *Store) GetSlice(ctx context.Context, key string) ([]string, error) {
	saveKey := fmt.Sprintf("%s:%s", c.keyPrefix, key)
	result, err := c.client.LRange(ctx, saveKey, 0, -1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, apierr.InternalServer().WithMsg("redis get slice failed").WithErr(err).WithStack()
	}
	return result, nil
}

func (c *Store) SetSlice(ctx context.Context, key string, value []any, expireTime *time.Duration) error {
	saveKey := fmt.Sprintf("%s:%s", c.keyPrefix, key)
	if expireTime == nil {
		if err := c.client.RPush(ctx, saveKey, value...).Err(); err != nil {
			return apierr.InternalServer().WithStack().WithMsg(fmt.Sprintf("redis setting %s key failed", saveKey)).WithErr(err)
		}
		return nil
	}
	if expireTime == &NeverExpires {
		if err := c.client.RPush(ctx, saveKey, value...).Err(); err != nil {
			return apierr.InternalServer().WithStack().WithMsg(fmt.Sprintf("redis setting %s key failed", saveKey)).WithErr(err)
		}
		return nil
	}
	return nil
}

func (c *Store) GetString(ctx context.Context, key string) (string, error) {
	saveKey := fmt.Sprintf("%s:%s", c.keyPrefix, key)
	if v, err := c.client.Get(ctx, saveKey).Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", apierr.InternalServer().WithStack().WithMsg(fmt.Sprintf("redis getting %s key failed", saveKey)).WithErr(err)
	} else {
		return v, nil
	}
}

// SetString 设置字符串
//
// expireTime 过期时间, nil 使用默认过期时间; &data.NeverExpires 表示永不过期
func (c *Store) SetString(ctx context.Context, key string, value string, expireTime *time.Duration) error {
	saveKey := fmt.Sprintf("%s:%s", c.keyPrefix, key)
	if expireTime == nil {
		if err := c.client.Set(ctx, saveKey, value, c.expireTime).Err(); err != nil {
			return apierr.InternalServer().WithStack().WithMsg(fmt.Sprintf("redis setting %s key failed", saveKey)).WithErr(err)
		}
		return nil
	}
	if expireTime == &NeverExpires {
		if err := c.client.Set(ctx, saveKey, value, 0).Err(); err != nil {
			return apierr.InternalServer().WithStack().WithMsg(fmt.Sprintf("redis setting %s key failed", saveKey)).WithErr(err)
		}
		return nil
	}
	if err := c.client.Set(ctx, saveKey, value, *expireTime).Err(); err != nil {
		if err = c.client.Set(ctx, saveKey, value, 0).Err(); err != nil {
			return apierr.InternalServer().WithStack().WithMsg(fmt.Sprintf("redis setting %s key failed", saveKey)).WithErr(err)
		}
		return nil
	}
	return nil
}

func (c *Store) GetInt64(ctx context.Context, key string) (*int64, error) {
	saveKey := fmt.Sprintf("%s:%s", c.keyPrefix, key)
	v, err := c.client.Get(ctx, saveKey).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, apierr.InternalServer().WithStack().WithMsg(fmt.Sprintf("redis getting %s key failed", saveKey)).WithErr(err)
	}
	return &v, nil
}

// SetInt64 设置整数
//
// expireTime 过期时间, nil 使用默认过期时间; &data.NeverExpires 表示永不过期
func (c *Store) SetInt64(ctx context.Context, key string, value int64, expireTime *time.Duration) error {
	saveKey := fmt.Sprintf("%s:%s", c.keyPrefix, key)
	if expireTime == nil {
		if err := c.client.Set(ctx, saveKey, value, c.expireTime).Err(); err != nil {
			return apierr.InternalServer().WithStack().WithMsg(fmt.Sprintf("redis setting %s key failed", saveKey)).WithErr(err)
		}
		return nil
	}
	if expireTime == &NeverExpires {
		if err := c.client.Set(ctx, saveKey, value, 0).Err(); err != nil {
			return apierr.InternalServer().WithStack().WithMsg(fmt.Sprintf("redis setting %s key failed", saveKey)).WithErr(err)
		}
		return nil
	}
	if err := c.client.Set(ctx, saveKey, value, *expireTime).Err(); err != nil {
		return apierr.InternalServer().WithStack().WithMsg(fmt.Sprintf("redis setting %s key failed", saveKey)).WithErr(err)
	}
	return nil
}

func (c *Store) Del(ctx context.Context, key string) error {
	saveKey := fmt.Sprintf("%s:%s", c.keyPrefix, key)
	if err := c.client.Del(ctx, saveKey).Err(); err != nil {
		return apierr.InternalServer().WithStack().WithMsg(fmt.Sprintf("redis deleting %s key failed", saveKey)).WithErr(err)
	}
	return nil
}

func (c *Store) Flush(ctx context.Context) error {
	if err := c.client.FlushDB(ctx).Err(); err != nil {
		return apierr.InternalServer().WithStack().WithMsg("redis flush db failed").WithErr(err)
	}
	return nil
}

func (c *Store) Incr(ctx context.Context, key string) (int64, error) {
	saveKey := fmt.Sprintf("%s:%s", c.keyPrefix, key)
	return c.client.Incr(ctx, saveKey).Result()
}
