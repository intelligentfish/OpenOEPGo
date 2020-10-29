package user

// User
type User struct {
	Id   int    `json:"id"`
	Role uint32 `json:"role"`
}

// IsRoot
func (object *User) IsRoot() bool {
	return IsRole(object.Role, RoleManager)
}

// CanPublishVideo
func (object *User) CanPublishVideo() bool {
	return IsRole(object.Role, RoleTeacher)
}

// CanSubscribeVideo
func (object *User) CanSubscribeVideo() bool {
	return IsRole(object.Role, RoleTeacher) || IsRole(object.Role, RoleStudent)
}
