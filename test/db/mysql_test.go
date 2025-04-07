package db

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"log"
	"qqlx/base/conf"
	"qqlx/base/data"
	"qqlx/base/logger"
	"qqlx/model"
	"qqlx/store/rbac"
	"qqlx/store/userstore"
	"testing"
)

var (
	f   func()
	err error
	sql *gorm.DB
)

func InitMysql() {
	err = conf.LoadConfig("../../config.yaml")
	if err != nil {
		log.Fatalf("load config faild: %v", err)
	}
	logger.InitLogger()
	sql, f, err = data.InitMySQL()
	if err != nil {
		log.Fatalf("init mysql faild: %v", err)
	}
}

func TestCreateTable(t *testing.T) {
	InitMysql()
	err = sql.AutoMigrate(&model.User{}, &model.Role{}, &model.Policy{})
	if err != nil {
		t.Fatalf("create table faild: %v", err)
	}
	defer f()
}

func TestRoleLoadUsers(t *testing.T) {
	InitMysql()
	store := rbac.NewRoleStore(sql)
	query, err := store.Query(context.Background(), rbac.LoadPolices(), rbac.LoadUsers(), rbac.RoleID(939524100))
	if err != nil {
		return
	}
	fmt.Printf("query: %#v", query)
}

func TestUserLoadRole(t *testing.T) {
	InitMysql()
	store := userstore.NewUserStore(sql)
	query, err := store.Query(context.Background(), userstore.LoadRoles(), userstore.ID(1543503876))
	if err != nil {
		return
	}
	fmt.Printf("query: %#v", query)
}
