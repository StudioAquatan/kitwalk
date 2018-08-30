package kitwalk

import "fmt"

// InvalidUsernameError will be return when given user name is invalid.
type InvalidUsernameError struct {
	username string
}

func (e *InvalidUsernameError) Error() string {
	return fmt.Sprintf(`Given username '%s' is invalid. For the user name, 
			combine a prefix such as 'b' or 'm' with the student number 
			excluding the first digit.`, e.username)
}

// ConfigDoesNotExists will raise when some configurations are missing.
type ConfigDoesNotExists struct{}

func (e *ConfigDoesNotExists) Error() string {
	return "Several configuration parameters do not exist."
}

// ShibbolethAuthError will raise when authentication has been failed.
type ShibbolethAuthError struct {
	username string
	errMsg   string
}

func (e *ShibbolethAuthError) Error() string {
	return fmt.Sprintf(`Authentication failed. 
			Given error message from auth page is as follows.\n%+v`, e.errMsg)
}
