package kitwalk

import (
	"net/url"
)

const (
	// DefaultUnameKey is used to post auth information
	DefaultUnameKey = "j_username"
	// DefaultPasswdKey is used as well as DefaultUnameKey
	DefaultPasswdKey = "j_password"
	// DefaultRelayStateKey is the key to parse HTML and extract saml auth information.
	DefaultRelayStateKey = "RelayState"
	// DefaultSAMLResponseKey is used as well as DefaultRelayStateKey
	DefaultSAMLResponseKey = "SAMLResponse"
	// DefaultAuthDomain is the domain of auth server.
	DefaultAuthDomain = "auth.cis.kit.ac.jp"
	// ShibbolethLoginURL is the default login url.
	ShibbolethLoginURL = "https://portal.student.kit.ac.jp/"
	// This pair will send with username and password.
	// The auth server will require this params.
	eventIDProceedKey = "_eventId_proceed"
	eventIDProceedVal = ""
	// These values are required to skip confirmation page.
	shibIdpLsExceptionKey = "shib_idp_ls_exception.shib_idp_session_ss"
	shibIdpLsExceptionVal = ""
	shibIdpLsSuccessKey   = "shib_idp_ls_success.shib_idp_session_ss"
	shibIdpLsSuccessVal   = "true"
)

// Config struct will have settings for saml authentication.
type Config struct {
	// Username key is used when this module POST auth information to auth server.
	ShibbolethUsernameKey string
	// Password key is used as well as Username key.
	ShibbolethPasswordKey string
	// Domain information of auth server.
	ShibbolethAuthDomain string
	// Url to login
	ShibbolethLoginURL string
	// This params has an additionally information to auth.
	// POST with username and password.
	ShibbolethHiddenParams url.Values
	// When appear webstorage confirmation during authentication steps, this params send to the server.
	ShibbolethPassConfirmationParams url.Values
}

// GetDefaultConfig will return the default configuration. It is enough to authenticate typically.
func GetDefaultConfig() *Config {
	var (
		defaultHiddenParams           = url.Values{}
		defaultPassConfirmationParams = url.Values{}
	)
	defaultHiddenParams.Add(eventIDProceedKey, eventIDProceedVal)
	defaultPassConfirmationParams.Add(shibIdpLsExceptionKey, shibIdpLsExceptionVal)
	defaultPassConfirmationParams.Add(shibIdpLsSuccessKey, shibIdpLsSuccessVal)
	defaultPassConfirmationParams.Add(eventIDProceedKey, eventIDProceedVal)
	return &Config{
		ShibbolethUsernameKey:            DefaultUnameKey,
		ShibbolethPasswordKey:            DefaultPasswdKey,
		ShibbolethAuthDomain:             DefaultAuthDomain,
		ShibbolethLoginURL:               ShibbolethLoginURL,
		ShibbolethHiddenParams:           defaultHiddenParams,
		ShibbolethPassConfirmationParams: defaultPassConfirmationParams,
	}
}
