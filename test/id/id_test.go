package id_test

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/sony/sonyflake"
	"log"
	"math/rand"
	"net"
	"os"
	"qqlx/base/conf"
	"qqlx/base/data"
	"qqlx/store/cache"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestGenId(t *testing.T) {
	s := sonyflake.Settings{
		StartTime: time.Now(),
		MachineID: func() (uint16, error) {
			return getLocalIPLow16Bits()
		},
	}
	so, err := sonyflake.New(s)
	if err != nil {
		t.Fatalf("初始化sonyflake失败: %s", err.Error())
	}
	id, err := so.NextID()
	if err != nil {
		t.Fatalf("生成id失败: %s", err.Error())
	}
	fmt.Printf("ulid: %v\n", id)
	id, err = so.NextID()
	if err != nil {
		t.Fatalf("生成id失败: %s", err.Error())
	}
	fmt.Printf("ulid: %v", id)

}

func getLocalIPLow16Bits() (uint16, error) {
	// 获取所有网络接口的地址
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return 0, err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ip := ipNet.IP.To4()
			ipInt := binary.BigEndian.Uint32(ip) // 将 IP 转换为 uint32
			low16 := uint16(ipInt & 0xFFFF)      // 取低16位
			return low16, nil
		}
	}
	return 0, fmt.Errorf("未找到可用的 IPv4 地址")
}

func TestUUID(t *testing.T) {
	uid, _ := uuid.NewUUID()
	fmt.Printf("uid: %s", uid)
}

func generateULID() string {
	// 使用时间戳作为种子
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id := ulid.MustNew(ulid.Now(), entropy)
	return strings.ToLower(id.String())
}

func TestUID(t *testing.T) {
	id := generateULID()
	fmt.Println(id)
}

func getMachineID() (uint16, error) {
	ctx := context.Background()
	rdb, err := data.CreateRDB(ctx)
	if err != nil {
		log.Fatalf("init redis faild: %v", err)
	}
	redisCli, f, err := cache.NewStore(rdb)
	if err != nil {
		log.Fatalf("init redis faild: %v", err)
	}
	defer f()

	id, err := redisCli.Incr(ctx, "machine_id")
	if err != nil {
		return 0, err
	}

	if id > 65535 {
		return 0, fmt.Errorf("machine ID overflow")
	}

	return uint16(id), nil
}

func TestGenID(t *testing.T) {
	err := conf.LoadConfig("../../config.yaml")
	if err != nil {
		t.Fatalf("load config faild: %v", err)
	}
	settings := sonyflake.Settings{
		StartTime: time.Now(),
		MachineID: getMachineID,
	}
	var result sync.Map

	sf := sonyflake.NewSonyflake(settings)
	wg := &sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			id, err := sf.NextID()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			result.Store(id, struct{}{})
			wg.Done()
		}()
	}

	wg.Wait()
	count := 0
	result.Store(1, struct{}{})
	result.Store(1, struct{}{})
	result.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	fmt.Printf("len: %d\n", count)
}
