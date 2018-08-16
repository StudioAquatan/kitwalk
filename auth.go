package kitwalk

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

type Auth interface {
	LoginWith(client *http.Client) error
	SetupWith(config Config) error
	LoginAs(username string, password string) error
}

type User struct {
	Username string
	Password string
}

type SamlAuthenticator struct {
	User   *User
	Config Config
}

func (c *SamlAuthenticator) auth(client *http.Client, resp *http.Response) error {
	var (
		err error
	)
	defer resp.Body.Close()
	// When Web Storage confirmation page appear, skip it.
	if isContinueRequired(resp.Body) {
		resp, err = client.PostForm(resp.Request.URL.String(), c.Config.ShibbolethPassConfirmationParams)
		if err != nil {
			return err
		}
		defer func() {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}()
	}

	// Post auth info to auth page
	authResp, err := client.PostForm(resp.Request.URL.String(), c.Config.ShibbolethHiddenParams)
	if err != nil {
		return err
	}
	defer authResp.Body.Close()
	// Extract SAML response
	actionUrl, data, err := parseSamlResp(authResp.Body)
	if err != nil {
		return err
	}
	// Redirect to target resource, and respond with target resource.
	authResult, err := client.PostForm(actionUrl, data)
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, authResult.Body)
		authResult.Body.Close()
	}()
	if authResult.Request.URL.Host == DefaultAuthDomain {
		return &ShibbolethAuthError{errMsg: "Try to auth, but return login page yet."}
	}
	return nil
}

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
	resp, err := client.Get(ShibbolethLoginUrl)
	if err != nil {
		return err
	}
	if resp.Request.URL.Host == c.Config.ShibbolethAuthDomain {
		return c.auth(client, resp)
	}
	return nil
}

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
