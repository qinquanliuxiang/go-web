package data_test

import (
	"qqlx/model"
	"testing"
)

func TestCreateTable(t *testing.T) {
	db.AutoMigrate(&model.User{}, &model.Role{}, &model.Policy{})
}
