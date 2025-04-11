package conf_test

import (
	"qqlx/base/conf"
	"testing"
)

func TestViper(t *testing.T) {
	err := conf.LoadConfig("../../config.yaml")
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}
	d := conf.GetMysqlMaxLifetime()
	s := d.String()
	t.Logf("time string: %v", s)
}
