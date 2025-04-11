package sonyflake

import (
	"context"
	"fmt"
	"github.com/sony/sonyflake"
	"qqlx/base/apierr"
	"qqlx/base/constant"
	"qqlx/store/cache"
	"time"
)

type GenerateIDStruct struct {
	sonyflake *sonyflake.Sonyflake
}

func NewGenerateID(ctx context.Context, redis *cache.Store) *GenerateIDStruct {
	return &GenerateIDStruct{
		sonyflake: sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime: time.Now(),
			MachineID: func() (uint16, error) {
				id, err := redis.Incr(ctx, constant.DefaultRedisIncrKey)
				if err != nil {
					return 0, fmt.Errorf("incr redis failed: %v", err)
				}
				return uint16(id), nil
			},
			CheckMachineID: func(id uint16) bool {
				if id == 0 || id > 65535 {
					return false
				}
				return true
			},
		}),
	}

}

func (g *GenerateIDStruct) NextID() (int, error) {
	id, err := g.sonyflake.NextID()
	if err != nil {
		return 0, apierr.InternalServer().Set(apierr.SonyflakeErrCode, "sonyflake next id failed", err)
	}
	return int(id), nil
}
