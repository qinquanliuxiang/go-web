package store

import (
	"qqlx/base/data"
	"qqlx/base/interfaces"
	"qqlx/pkg/sonyflake"
	"qqlx/store/cache"
	"qqlx/store/ldap"
	"qqlx/store/rbac"
	"qqlx/store/userstore"

	"github.com/google/wire"
)

var ProviderStore = wire.NewSet(
	wire.Bind(new(interfaces.CacheInterface), new(*cache.Store)),
	wire.Bind(new(interfaces.UserStoreInterface), new(*userstore.Store)),
	wire.Bind(new(interfaces.UserRoleStoreInterface), new(*userstore.UserAssociationStore)),
	wire.Bind(new(interfaces.RoleStoreInterface), new(*rbac.RoleStore)),
	wire.Bind(new(interfaces.PolicyStoreInterface), new(*rbac.PolicyStore)),
	wire.Bind(new(interfaces.RolePolicyStoreInterface), new(*rbac.RoleAssociationStore)),
	wire.Bind(new(interfaces.CasbinInterface), new(*rbac.CasbinStore)),
	wire.Bind(new(interfaces.LdapInterface), new(*ldap.Store)),
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
