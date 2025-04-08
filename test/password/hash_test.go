package password_test

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"qqlx/base/conf"
	"qqlx/base/constant"
	"qqlx/base/data"
	"qqlx/base/logger"
	"qqlx/pkg/sonyflake"
	"qqlx/schema"
	"qqlx/service"
	"qqlx/store/cache"
	"qqlx/store/ldap"
	"qqlx/store/userstore"
	"testing"

	"github.com/google/uuid"
)

// 对密码进行加盐哈希
func hashPassword(password string, salt []byte) string {
	// 将盐和密码混合在一起
	combined := append([]byte(password), salt...)
	// 使用 SHA-256 搅拌成最终的哈希值
	hash := sha256.Sum256(combined)
	return hex.EncodeToString(hash[:])
}

func TestSha256Crypt(t *testing.T) {
	// 用户输入的密码
	password := "123456"
	_slat := "r7gLpyAu"
	slat := []byte(_slat)
	// 对密码进行加盐哈希
	hashedPassword := hashPassword(password, slat)

	fmt.Println("Password:", password)
	fmt.Println("Salt:", hex.EncodeToString(slat))
	fmt.Println("Hashed Password:", base64.StdEncoding.EncodeToString([]byte(hashedPassword)))
}

func TestLdapSha256Crypt(t *testing.T) {
	err := conf.LoadConfig("../../config.yaml")
	if err != nil {
		t.Fatalf("load config faild: %v", err)
	}
	logger.InitLogger()
	ldapCon, f1, err := data.InitLdap()
	if err != nil {
		t.Fatalf("init ldap faild: %v", err)
	}
	mysql, f2, err := data.InitMySQL()
	if err != nil {
		t.Fatalf("init mysql faild: %v", err)
	}
	defer func() {
		f1()
		f2()
	}()
	userStore := userstore.NewUserStore(mysql)
	ldapStore, err := ldap.NewLdapStore(ldapCon)
	if err != nil {
		t.Fatalf("new ldap store faild: %v", err)
	}
	redisCli, err := data.CreateRDB(context.Background())
	if err != nil {
		t.Fatalf("init redis faild: %v", err)
	}
	cacheStore, f3, err := cache.NewStore(redisCli)
	if err != nil {
		t.Fatalf("init cache store faild: %v", err)
	}
	defer f3()
	generateID := sonyflake.NewGenerateID(context.Background(), cacheStore)
	userSVC, err := service.NewUserSVC(generateID, userStore, nil, nil, cacheStore, nil, ldapStore)
	if err != nil {
		t.Fatalf("new user svc faild: %v", err)
	}

	req := &schema.UserRegistryRequest{
		Name:     "test1",
		Password: "123456",
		Avatar:   "avatar",
		Email:    "test1@qqlx.net",
		Mobile:   "mobile",
	}
	requestID := uuid.New().String()
	ctx := context.WithValue(context.TODO(), constant.TraceID, requestID)
	err = userSVC.RegistryUser(ctx, req)
	if err != nil {
		t.Fatalf("registry user faild: %v", err)
	}
}

// 生成 SSHA 加密的密码
func encryptSSHA(password string) (string, error) {
	// 生成盐值
	salt := "123456"
	// 计算 SHA1 哈希
	hash := sha1.New()
	hash.Write([]byte(password))
	hash.Write([]byte(salt))

	// 获取 SHA1 哈希值
	hashResult := hash.Sum(nil)

	// 将哈希值和盐值组合
	result := append(hashResult, salt...)

	// 将结果进行 Base64 编码
	encoded := base64.StdEncoding.EncodeToString(result)

	// 返回 SSHA 格式的密码
	return "{SSHA}" + encoded, nil
}

func TestEncryptSSHA(t *testing.T) {
	password := "123456"
	encryptedPassword, err := encryptSSHA(password)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}
	fmt.Println("加密后的密码:", encryptedPassword)
}
