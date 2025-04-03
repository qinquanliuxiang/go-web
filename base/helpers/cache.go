package helpers

import (
	"fmt"
	"qqlx/base/constant"
)

func GetRoleCacheKey(name string) string {
	return fmt.Sprintf("%s:%s", constant.RoleCacheKeyPrefix, name)
}
