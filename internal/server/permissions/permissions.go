package permissions

type Perms []Perm

//go:generate go run golang.org/x/tools/cmd/stringer -type=Perm -trimprefix=Perm
type Perm int32
const (
	PermUsersDelete                Perm = 1
	PermUsersUpdateName            Perm = 2
	PermUsersUpdateNameOwn         Perm = 3
	PermUsersUpdateRole            Perm = 4
	PermRolesCreate                Perm = 5
	PermRolesDelete                Perm = 6
	PermRolesUpdateName            Perm = 7
	PermRolesUpdatePerms           Perm = 8
)

var AllPerms Perms = Perms{
	PermUsersDelete,
	PermUsersUpdateName,
	PermUsersUpdateNameOwn,
	PermUsersUpdateRole,
	PermRolesCreate,
	PermRolesDelete,
	PermRolesUpdateName,
	PermRolesUpdatePerms,
}
