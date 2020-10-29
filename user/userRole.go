package user

const (
	RoleNormal  = uint32(0x100) // normal role
	RoleStudent = uint32(0x200) // student role
	RoleTeacher = uint32(0x300) // teacher role
	RoleManager = uint32(0x400) // manager role
)

// IsRole
func IsRole(value, role uint32) bool {
	return role == (value & role)
}
