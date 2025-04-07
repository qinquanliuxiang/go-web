package store

import (
	"qqlx/base/data"
	"qqlx/pkg/sonyflake"
	"qqlx/store/cache"
	"qqlx/store/ldap"
	"qqlx/store/rbac"
	"qqlx/store/userstore"

	"github.com/google/wire"
)

var ProviderStore = wire.NewSet(
	wire.Bind(new(CacheInterface), new(*cache.Store)),
	wire.Bind(new(UserStoreInterface), new(*userstore.Store)),
	wire.Bind(new(UserRoleStoreInterface), new(*userstore.UserAssociationStore)),
	wire.Bind(new(RoleStoreInterface), new(*rbac.RoleStore)),
	wire.Bind(new(PolicyStoreInterface), new(*rbac.PolicyStore)),
	wire.Bind(new(RolePolicyStoreInterface), new(*rbac.RoleAssociationStore)),
	wire.Bind(new(CasbinInterface), new(*rbac.CasbinStore)),
	wire.Bind(new(LdapInterface), new(*ldap.Store)),
	data.CreateRDB,
	data.InitMySQL,
	data.InitLdap,
	cache.NewStore,
	userstore.NewUserStore,
	userstore.NewUserAssociationStore,
	rbac.NewRoleStore,
	rbac.NewPolicyStore,
	rbac.NewRoleAssociationStore,
	ldap.NewLdapStore,
	rbac.NewCasbinStore,
	data.InitCasbin,
	sonyflake.NewGenerateID,
)
