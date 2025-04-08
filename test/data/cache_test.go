package data_test

import (
	"context"
	"qqlx/store/cache"
	"testing"
)

func TestCacge(t *testing.T) {
	// defer test.Close1()
	// defer test.Close2()
	if err := cacheImpl.SetInt64(context.Background(), "test", 100, &cache.NeverExpires); err != nil {
		t.Logf("set err: %v", err)
	}
	defer f()
	// if err := test.Cache.Del(context.Background(), "test"); err != nil {
	// 	t.Logf("del err: %v", err)
	// }
}

func TestCacge2(t *testing.T) {
	InitCli()
	defer f()
	err := cacheImpl.SetSet(ctx, "test", []any{1, 2, 3}, nil)
	if err != nil {
		t.Fatalf("set err: %v", err)
	}
	t.Logf("set success")

	err = cacheImpl.SetSet(ctx, "test", []any{1, 2, 3}, nil)
	if err != nil {
		t.Fatalf("set err: %v", err)
	}
	t.Logf("set success")

	slice, err := cacheImpl.GetSet(ctx, "test")
	if err != nil {
		t.Fatalf("get err: %v", err)
	}
	t.Logf("get success: %v", slice)
}
