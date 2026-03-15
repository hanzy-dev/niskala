package auth

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

func IsAdmin(role string) bool {
	return role == RoleAdmin
}
