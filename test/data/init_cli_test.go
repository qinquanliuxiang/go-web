package data_test

import (
	"context"
	"qqlx/base/conf"
	"qqlx/base/data"
	"qqlx/store"
	"qqlx/store/cache"

	"gorm.io/gorm"
)

var (
	ctx       = context.Background()
	db        *gorm.DB
	f         func()
	cacheImpl store.CacheInterface
)

func InitCli() {
	err := conf.LoadConfig("../../config.yaml")
	if err != nil {
		panic(err)
	}
	mysqlCli, close1, err := data.InitMySQL()
	if err != nil {
		panic(err)
	}
	rdb, err := data.CreateRDB(ctx)
	if err != nil {
		panic(err)
	}
	redisCli, close2, err := cache.NewStore(rdb)
	if err != nil {
		panic(err)
	}
	f = func() {
		close1()
		close2()
	}
	db = mysqlCli
	cacheImpl = redisCli
}
