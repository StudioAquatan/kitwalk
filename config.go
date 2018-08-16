package kitwalk

import (
	"net/url"
)

const (
	DefaultUnameKey        = "j_username"
	DefaultPasswdKey       = "j_password"
	DefaultRelayStateKey   = "RelayState"
	DefaultSAMLResponseKey = "SAMLResponse"
	DefaultAuthDomain      = "auth.cis.kit.ac.jp"
	EventIdProceedKey      = "_eventId_proceed"
	EventIdProceedVal      = ""
	ShibbolethLoginUrl     = "https://portal.student.kit.ac.jp/"
	ShibIdpLsExceptionKey  = "shib_idp_ls_exception.shib_idp_session_ss"
	ShibIdpLsExceptionVal  = ""
	ShibIdpLsSuccessKey    = "shib_idp_ls_success.shib_idp_session_ss"
	ShibIdpLsSuccessVal    = "true"
)

type Config struct {
	ShibbolethUsernameKey            string
	ShibbolethPasswordKey            string
	ShibbolethAuthDomain             string
	ShibbolethLoginUrl               string
	ShibbolethHiddenParams           url.Values
	ShibbolethPassConfirmationParams url.Values
}

func GetDefaultConfig() *Config {
	var (
		defaultHiddenParams           = url.Values{}
		defaultPassConfirmationParams = url.Values{}
	)
	defaultHiddenParams.Add(EventIdProceedKey, EventIdProceedVal)
	defaultPassConfirmationParams.Add(ShibIdpLsExceptionKey, ShibIdpLsExceptionVal)
	defaultPassConfirmationParams.Add(ShibIdpLsSuccessKey, ShibIdpLsSuccessVal)
	defaultPassConfirmationParams.Add(EventIdProceedKey, EventIdProceedVal)
	return &Config{
		ShibbolethUsernameKey:            DefaultUnameKey,
		ShibbolethPasswordKey:            DefaultPasswdKey,
		ShibbolethAuthDomain:             DefaultAuthDomain,
		ShibbolethLoginUrl:               ShibbolethLoginUrl,
		ShibbolethHiddenParams:           defaultHiddenParams,
		ShibbolethPassConfirmationParams: defaultPassConfirmationParams,
	}
}
