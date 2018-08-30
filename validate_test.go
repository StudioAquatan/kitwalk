package kitwalk

import "testing"

func TestIsValidUsername(t *testing.T) {
	t.Parallel()
	var (
		prefixes    = [...]string{"b", "m", "d"}
		num         = "1234567"
		numTooShort = "123456"
		numTooLong  = "12345678"
	)
	t.Run("Valid username test", func(t *testing.T) {
		for _, prefix := range prefixes {
			uName := prefix + num
			err := isValidUsername(uName)
			if err != nil {
				t.Errorf("It recognized '%s' as invalid", uName)
			}
		}
	})
	t.Run("Too short username test", func(t *testing.T) {
		for _, prefix := range prefixes {
			uName := prefix + numTooShort
			err := isValidUsername(uName)
			if err == nil {
				t.Errorf("It recognized '%s' as valid", uName)
			}
		}
	})
	t.Run("Too long username test", func(t *testing.T) {
		for _, prefix := range prefixes {
			uName := prefix + numTooLong
			err := isValidUsername(uName)
			if err == nil {
				t.Errorf("It recognized '%s' as valid", uName)
			}
		}
	})
}
