package data_test

import (
	"qqlx/model"
	"testing"
)

func TestCreateTable(t *testing.T) {
	InitCli()
	defer f()
	db.AutoMigrate(&model.User{}, &model.Role{}, &model.Policy{})
}
