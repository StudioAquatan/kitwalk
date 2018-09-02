package kitwalk

import (
	"net/http"
	"net/http/cookiejar"
	"strings"
)

const (
	contentTypeHead = "Content-Type"
	contentTypeVal  = "application/x-www-form-urlencoded"
)

// Auth is an interface for http.Client
type Auth interface {
	LoginWith(client *http.Client) error
	SetupWith(config Config) error
	LoginAs(username string, password string) error
}

// User is an user belonging to the authentication destination
type User struct {
	Username string
	Password string
}

// SamlAuthenticator has Config and User. This struct implement Auth interface.
type SamlAuthenticator struct {
	User   *User
	Config Config
}

func (c *SamlAuthenticator) auth(client *http.Client, resp *http.Response) error {
	var (
		err     error
		tmpResp = resp
	)
	// When Web Storage confirmation page appear, skip it.
	if isContinueRequired(resp.Body) {
		param := strings.NewReader(c.Config.ShibbolethPassConfirmationParams.Encode())
		crReq, err := http.NewRequest(http.MethodPost, tmpResp.Request.URL.String(), param)
		if err != nil {
			return err
		}
		crReq.Header.Add(contentTypeHead, contentTypeVal)
		tmpResp, err = client.Do(crReq)
		if err != nil {
			return err
		}
		defer tmpResp.Body.Close()
	}

	// Post auth info to auth page
	param := strings.NewReader(c.Config.ShibbolethHiddenParams.Encode())
	authReq, err := http.NewRequest(http.MethodPost, tmpResp.Request.URL.String(), param)
	if err != nil {
		return err
	}
	authReq.Header.Add(contentTypeHead, contentTypeVal)
	authResp, err := client.Do(authReq)
	if err != nil {
		return err
	}
	defer authResp.Body.Close()
	// Extract SAML response
	actionURL, data, err := parseSamlResp(authResp.Body)
	if err != nil {
		return err
	}
	// Redirect to target resource, and respond with target resource.
	param = strings.NewReader(data.Encode())
	authResReq, err := http.NewRequest(http.MethodPost, actionURL, param)
	if err != nil {
		return err
	}
	authResReq.Header.Add(contentTypeHead, contentTypeVal)
	authResult, err := client.Do(authResReq)
	if err != nil {
		return err
	}
	defer authResult.Body.Close()
	if authResult.Request.URL.Host == DefaultAuthDomain {
		return &ShibbolethAuthError{errMsg: "Try to auth, but return login page yet."}
	}
	return nil
}

// LoginWith works with given http.Client to auth.
// The client store cookie information to be used for next authentication.
func (c *SamlAuthenticator) LoginWith(client *http.Client) error {
	if client == nil {
		client = http.DefaultClient
	}
	if client.Jar == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return err
		}
		client.Jar = jar
	}
	getReq, err := http.NewRequest(http.MethodGet, ShibbolethLoginURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(getReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Request.URL.Host == c.Config.ShibbolethAuthDomain {
		err = c.auth(client, resp)
		return err
	}
	return nil
}

// SetupWith attach given configuration to authenticator.
func (c *SamlAuthenticator) SetupWith(config Config) error {
	// Set auth info
	if config.ShibbolethHiddenParams == nil && config.ShibbolethPassConfirmationParams == nil {
		return &ConfigDoesNotExists{}
	}
	authParams := config.ShibbolethHiddenParams
	if len(authParams.Get(config.ShibbolethUsernameKey)) == 0 {
		authParams.Add(config.ShibbolethUsernameKey, c.User.Username)
	} else {
		authParams.Set(config.ShibbolethUsernameKey, c.User.Username)
	}
	if len(authParams.Get(config.ShibbolethPasswordKey)) == 0 {
		authParams.Add(config.ShibbolethPasswordKey, c.User.Password)
	} else {
		authParams.Set(config.ShibbolethPasswordKey, c.User.Password)
	}
	c.Config = config
	return nil
}

// LoginAs switch user to authenticate with
func (c *SamlAuthenticator) LoginAs(username string, password string) error {
	if err := isValidUsername(username); err != nil {
		return err
	}
	user := &User{Username: username, Password: password}
	c.User = user
	err := c.SetupWith(c.Config)
	if err != nil {
		return err
	}
	return nil
}

// NewAuthenticator create new authenticator with given auth information.
func NewAuthenticator(username string, password string) (Auth, error) {
	if err := isValidUsername(username); err != nil {
		return nil, err
	}
	user := &User{Username: username, Password: password}
	defaultConfig := GetDefaultConfig()
	authenticator := &SamlAuthenticator{
		User:   user,
		Config: *defaultConfig,
	}
	err := authenticator.SetupWith(*defaultConfig)
	if err != nil {
		return nil, err
	}
	return authenticator, nil
}
