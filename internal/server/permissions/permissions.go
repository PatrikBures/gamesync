package permissions

type Perms []Perm

// 1xxx users
// 2xxx roles

//go:generate stringer -type=Perm -trimprefix=Perm
type Perm int32
const (
	PermAllAllowed                  Perm = 10

	// /users
	PermUsersList					Perm = 1000
	// /users/{UserID}
	PermUserDelete                  Perm = 1050
	PermUserGet                     Perm = 1060 
	PermUserGetOwn                  Perm = 1070
	// /users/{UserID}/name
	PermUserNameUpdate              Perm = 1100
	PermUserNameUpdateOwn           Perm = 1110
	// /users/{UserID}/role
	PermUserRoleUpdate              Perm = 1130


	// /roles
	PermRolesGet                    Perm = 2000 // list roles, list role perms, get name
	PermRolesMod                    Perm = 2010 // mod perms/name for existing roles
	PermRolesDelete                 Perm = 2030 // delete role with no users
)

var AllPerms Perms = Perms{
	PermAllAllowed,
	PermUsersList,
	PermUserDelete,
	PermUserGet,
	PermUserGetOwn,
	PermUserNameUpdate,
	PermUserNameUpdateOwn,
	PermUserRoleUpdate,
	PermRolesGet,
	PermRolesMod,
	PermRolesDelete,
}
