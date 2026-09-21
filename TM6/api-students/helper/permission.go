package helper

// PermissionSet menyimpan daftar permission berdasarkan role.
type PermissionSet struct {
	permissions map[string]map[string]bool
}

// NewPermissionSet membuat PermissionSet dari data role dan permission.
func NewPermissionSet(rolePermissions map[string][]string) *PermissionSet {
	permissions := make(map[string]map[string]bool)

	for role, rolePerms := range rolePermissions {
		permissions[role] = make(map[string]bool)

		for _, permission := range rolePerms {
			permissions[role][permission] = true
		}
	}

	return &PermissionSet{
		permissions: permissions,
	}
}

// Has mengecek apakah suatu role memiliki permission tertentu.
func (p *PermissionSet) Has(role string, permission string) bool {
	if p == nil {
		return false
	}

	rolePermissions, ok := p.permissions[role]
	if !ok {
		return false
	}

	return rolePermissions[permission]
}
