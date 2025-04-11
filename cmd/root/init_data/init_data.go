package init_data

import (
	"context"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"os"
	"qqlx/base/conf"
	"qqlx/base/constant"
	"qqlx/base/data"
	"qqlx/base/logger"
	"qqlx/model"
	"qqlx/pkg/sonyflake"
	"qqlx/schema"
	"qqlx/service"
	"qqlx/store/cache"
	ldapstore "qqlx/store/ldap"
	"qqlx/store/rbac"
	"qqlx/store/userstore"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "init data",
	Long:  "init data",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed(constant.FlagConfigPath) {
			envConfigPath := os.Getenv(constant.ConfigEnv)
			if envConfigPath != "" {
				err := cmd.Flags().Set(constant.FlagConfigPath, envConfigPath)
				if err != nil {
					fmt.Printf("set config file path from env %s faild: %v", envConfigPath, err)
					return
				}
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cf, err := cmd.Flags().GetString(constant.FlagConfigPath)
		if err != nil {
			log.Fatalf("get config file path faild: %v", err)
		}
		initData(cf)
	},
}

func initData(cf string) {
	var (
		ldapCon   *ldap.Conn
		f         func()
		ldapStore *ldapstore.Store
		ctx       = context.Background()
	)
	ctxValue := context.WithValue(ctx, constant.TraceID, "init")
	err := conf.LoadConfig(cf)
	if err != nil {
		log.Fatalf("load config file %s failed: %v", cf, err)
	}
	ldapEnable := conf.GetLdapEnable()
	logger.InitLogger()
	db, closeFunc, err := data.InitMySQL()
	if err != nil {
		logger.Caller().Error("init mysql failed: %v", err)
		return
	}
	redisCli, err := data.CreateRDB(ctxValue)
	if err != nil {
		logger.Caller().Errorf("init redis faild: %v", err)
		return
	}
	enforcer, err := data.InitCasbin()
	if err != nil {
		logger.Caller().Errorf("init casbin faild: %v", err)
		return
	}

	if ldapEnable {
		ldapCon, f, err = data.InitLdap()
		if err != nil {
			logger.Caller().Errorf("init ldap faild: %v", err)
			return
		}
		defer f()
	}

	defer func() {
		_ = zap.S().Sync()
		closeFunc()
	}()
	if err = db.AutoMigrate(&model.User{}, &model.Role{}, &model.Policy{}); err != nil {
		panic(err)
	}
	casbinStore := rbac.NewCasbinStore(enforcer)
	userRepo := userstore.NewUserStore(db)
	userRoleStore := userstore.NewUserAssociationStore(db)
	roleRepo := rbac.NewRoleStore(db)
	cacheStore, f1, err := cache.NewStore(redisCli)
	defer f1()
	if err != nil {
		logger.Caller().Errorf("init cache store faild: %v", err)
	}
	if ldapEnable {
		ldapStore, err = ldapstore.NewLdapStore(ldapCon)
		if err != nil {
			logger.Caller().Errorf("init ldap store faild: %v", err)
		}
	}
	generateIDStruct := sonyflake.NewGenerateID(ctxValue, cacheStore)
	roleStore := rbac.NewRoleStore(db)
	policyStore := rbac.NewPolicyStore(db)
	appendStore := rbac.NewRoleAssociationStore(db)
	roleSvc := service.NewRoleSVC(generateIDStruct, roleStore, policyStore, appendStore, casbinStore, ldapStore)
	policySvc := service.NewPolicySVC(generateIDStruct, policyStore)
	// Create Polices
	for _, police := range polices {
		_ = policySvc.CreatePolicy(ctxValue, &police)
	}
	// Create Role
	err = roleSvc.CreateRole(ctxValue, &schema.RoleCreateRequest{
		Name:     "admin",
		Describe: "超级管理员",
	})
	if err != nil {
		logger.Caller().Error(err)
	}
	err = roleSvc.CreateRole(ctxValue, &schema.RoleCreateRequest{
		Name:     "view",
		Describe: "查看",
	})
	if err != nil {
		logger.Caller().Error(err)
	}

	adminRole, err := roleRepo.Query(ctxValue, rbac.RoleName("admin"))
	if err != nil {
		logger.Caller().Error(err)
	}
	adminPolicy, err := policyStore.Query(ctxValue, rbac.PolicyName("admin"))
	if err != nil {
		logger.Caller().Error(err)
	}

	// role 添加权限
	err = appendStore.AppendPolicy(ctxValue, adminRole, []model.Policy{*adminPolicy})
	if err != nil {
		logger.Caller().Error(err)
	}
	err = casbinStore.CreateRolePolices(ctxValue, [][]string{{adminRole.Name, adminPolicy.Path, adminPolicy.Method}})
	if err != nil {
		zap.S().Error(err)
	}

	viewRole, err := roleRepo.Query(ctxValue, rbac.RoleName("view"))
	if err != nil {
		logger.Caller().Error(err)
	}
	viewPolicy, err := policyStore.Query(ctxValue, rbac.PolicyName("view"))
	if err != nil {
		logger.Caller().Error(err)
	}

	err = appendStore.AppendPolicy(ctxValue, viewRole, []model.Policy{*viewPolicy})
	if err != nil {
		logger.Caller().Error(err)
	}
	err = casbinStore.CreateRolePolices(ctxValue, [][]string{{viewRole.Name, viewPolicy.Path, viewPolicy.Method}})
	if err != nil {
		zap.S().Error(err)
	}

	userSvc, err := service.NewUserSVC(generateIDStruct, userRepo, userRoleStore, roleRepo, cacheStore, casbinStore, ldapStore)
	if err != nil {
		logger.Caller().Error(err)
		return
	}
	_ = userSvc.RegistryUser(ctxValue, &schema.UserRegistryRequest{
		Name:     "admin",
		Password: "12345678",
		NickName: "超级管理员",
		Email:    "admin@example.com",
		Avatar:   "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Mobile:   "13800000000",
	})
	adminUser, err := userRepo.Query(ctxValue, userstore.Name("admin"))
	if err != nil {
		zap.S().Error(err)
		return
	}
	err = userSvc.UserAddRole(ctxValue, &schema.UserUpdateRoleRequest{ID: adminUser.ID, RoleNames: []string{"admin"}})
	if err != nil {
		zap.S().Error(err)
	}
}
