package auth

// CheckAuthAPI Check for auth - either for a (human) user, or an API user
func CheckAuthAPI() (bool, *User) {


	return false,nil
}

// CheckAuthUser Check for auth - endpoint only for (human) user use
func CheckAuthUser() (bool, *User) {
	return false,nil
}