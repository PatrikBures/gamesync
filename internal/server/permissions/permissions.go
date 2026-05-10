package permissions

type Perms []Perm

type Perm int
const (
	PermUsersDelete                Perm = 0
	PermUsersUpdateName            Perm = 1
	PermUsersUpdateNameOwn         Perm = 2
	PermUsersUpdateRole            Perm = 3
	PermRolesCreate                Perm = 4
	PermRolesDelete                Perm = 5
	PermRolesUpdateName            Perm = 6
	PermRolesUpdatePerms           Perm = 7
)
