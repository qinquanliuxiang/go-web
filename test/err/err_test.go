package err_test

import (
	"errors"
	"fmt"
	"testing"

	"gorm.io/gorm"
)

func TestErr(t *testing.T) {
	t.Run("err", func(t *testing.T) {
		err := gorm.ErrRecordNotFound
		newErr := fmt.Errorf("数据库错误: %w", err)
		t.Logf("newErr: %v", newErr)
		if errors.Is(newErr, gorm.ErrRecordNotFound) {
			t.Log("yes")
		} else {
			t.Log("no")
		}
	})
}
