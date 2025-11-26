package helper

func HasPermission(required string, perms []string) bool {
	for _, p := range perms {
		if p == required {
			return true
		}
	}
	return false
}

func HasRole(role string, required string) bool {
	return role == required
}