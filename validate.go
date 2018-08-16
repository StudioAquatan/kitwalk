package kitwalk

import "regexp"

var (
	usernameRegex = regexp.MustCompile("^[bmd]\\d{7}$")
)

func isValidUsername(username string) error {
	if !usernameRegex.MatchString(username) {
		return &InvalidUsernameError{username: username}
	}
	return nil
}
