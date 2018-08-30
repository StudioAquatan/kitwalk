package kitwalk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

const (
	invalidUsername = "testuser"
	invalidPasswd   = "invalidpasswd"
	validUsername   = "b1234567"
	validPasswd     = "passwd"
	validRelayState = "RelayState"
	validSAMLResp   = "SAMLResponse"
)

func check(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func TestNewAuthenticator(t *testing.T) {
	t.Run("Get new authenticator", func(t *testing.T) {
		_, err := NewAuthenticator(validUsername, validPasswd)
		check(t, err)
	})
	t.Run("Get authenticator with invalid username", func(t *testing.T) {
		_, err := NewAuthenticator(invalidUsername, validPasswd)
		switch e := err.(type) {
		case *InvalidUsernameError:
			// No problem
		default:
			t.Errorf("Expect error message is '%+v' but actual message is as follows \n '%+v'.", InvalidUsernameError{}, e)
		}
	})
}

type samlMock struct {
	Authenticated          bool
	WebStorageConfirmation bool
}

func (s *samlMock) RoundTrip(req *http.Request) (*http.Response, error) {
	// TODO: Refactoring
	resp := &http.Response{
		Header:  make(http.Header),
		Request: req,
	}
	if req.Method == http.MethodGet {
		if s.WebStorageConfirmation {
			wConf, err := ioutil.ReadFile("./samples/webstorage_confirm.html")
			if err != nil {
				return nil, err
			}
			req.URL, _ = url.Parse("https://auth.cis.kit.ac.jp/idp/profile/SAML2/Redirect/SSO?execution=e1s1")
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(wConf))
			return resp, nil
		}
		if !s.Authenticated {
			authForm, err := ioutil.ReadFile("./samples/auth_form.html")
			if err != nil {
				return nil, err
			}
			req.URL, _ = url.Parse("https://auth.cis.kit.ac.jp/idp/profile/SAML2/Redirect/SSO?execution=e1s1")
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(authForm))
			return resp, nil
		}
		portal, err := ioutil.ReadFile("./samples/internal_auth.html")
		if err != nil {
			return nil, err
		}
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(portal))
		return resp, nil
	} else if req.Method == http.MethodPost {
		if s.WebStorageConfirmation {
			err := req.ParseForm()
			if err != nil {
				return nil, err
			}
			q := req.PostForm
			successVal := q.Get(shibIdpLsSuccessKey)
			exceptionVal := q.Get(shibIdpLsExceptionKey)
			if successVal != shibIdpLsSuccessVal && exceptionVal != shibIdpLsExceptionVal {
				errTmpl := "\n[Expected]\n\t%+v: %+v, %+v: %+v\n\t%+v: %+v, %+v: %+v\n"
				errMsg := fmt.Sprintf(errTmpl, shibIdpLsSuccessKey, shibIdpLsSuccessVal, shibIdpLsExceptionKey, shibIdpLsExceptionVal,
					shibIdpLsSuccessKey, successVal, shibIdpLsExceptionKey, exceptionVal)
				return nil, errors.New(errMsg)
			}
			authForm, err := ioutil.ReadFile("./samples/auth_form.html")
			if err != nil {
				return nil, err
			}
			req.URL, _ = url.Parse("https://auth.cis.kit.ac.jp/idp/profile/SAML2/Redirect/SSO?execution=e1s1")
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(authForm))
			s.WebStorageConfirmation = false
			return resp, nil
		}
		if !s.Authenticated {
			err := req.ParseForm()
			if err != nil {
				return nil, err
			}
			q := req.PostForm
			uname := q.Get(DefaultUnameKey)
			passwd := q.Get(DefaultPasswdKey)
			if len(uname) == 0 && len(passwd) == 0 {
				relayState := q.Get(DefaultRelayStateKey)
				samlResponse := q.Get(DefaultSAMLResponseKey)
				if relayState == validRelayState && samlResponse == validSAMLResp {
					req.URL, _ = url.Parse(ShibbolethLoginURL)
					portal, err := ioutil.ReadFile("./samples/internal_auth.html")
					if err != nil {
						return nil, err
					}
					s.Authenticated = true
					resp.Body = ioutil.NopCloser(bytes.NewBuffer(portal))
					return resp, nil
				}
				errTmpl := "[Expected]\n\tRelayState: %s SAMLResponse: %s\n[Actual]\n\tRelayState: %s SAMLResponse: %s\n"
				errMsg := fmt.Sprintf(errTmpl, DefaultRelayStateKey, DefaultSAMLResponseKey, relayState, samlResponse)
				return nil, errors.New(errMsg)
			}
			if uname == validUsername && passwd == validPasswd {
				authSuccess, err := ioutil.ReadFile("./samples/auth_success.html")
				if err != nil {
					return nil, err
				}
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(authSuccess))
				return resp, nil
			}
			authFail, err := ioutil.ReadFile("./samples/auth_error.html")
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(authFail))
			return resp, nil
		}
		return resp, nil
	} else {
		return resp, nil
	}
}

func TestSamlAuthenticator_LoginWith(t *testing.T) {
	const (
		baseURL = "https://portal.student.kit.ac.jp/"
	)
	t.Run("Login with valid username and password", func(t *testing.T) {
		authenticator, err := NewAuthenticator(validUsername, validPasswd)
		check(t, err)
		client := &http.Client{
			Transport: &samlMock{Authenticated: false},
		}
		err = authenticator.LoginWith(client)
		check(t, err)
		resp, err := client.Get(baseURL)
		check(t, err)
		defer func() {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}()
		if resp.Request.URL.String() != baseURL {
			t.Errorf("Expect: [GET] %+v\nActual: %+v\n", resp.Request.URL, baseURL)
		}
	})
	t.Run("Login with invalid password", func(t *testing.T) {
		authenticator, err := NewAuthenticator(validUsername, invalidPasswd)
		check(t, err)
		client := &http.Client{
			Transport: &samlMock{Authenticated: false},
		}
		err = authenticator.LoginWith(client)
		switch e := err.(type) {
		case *ShibbolethAuthError:
			// Expected
		default:
			t.Errorf("Expected: ShibbolethAuthError\nActual: %+v\n", e)
		}
	})
	t.Run("Login with nil client", func(t *testing.T) {
		authenticator, err := NewAuthenticator(validUsername, validPasswd)
		check(t, err)
		client := http.DefaultClient
		client.Transport = &samlMock{Authenticated: false}
		err = authenticator.LoginWith(nil)
		check(t, err)
		resp, err := client.Get(baseURL)
		check(t, err)
		defer func() {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}()
		if resp.Request.URL.String() != baseURL {
			t.Errorf("Expect: [GET] %+v\nActual: %+v\n", resp.Request.URL, baseURL)
		}
	})
	t.Run("Skip WebStorage Confirmation", func(t *testing.T) {
		authenticator, err := NewAuthenticator(validUsername, validPasswd)
		check(t, err)
		client := &http.Client{
			Transport: &samlMock{
				Authenticated:          false,
				WebStorageConfirmation: true,
			},
		}
		err = authenticator.LoginWith(client)
		check(t, err)
	})
}

func TestSamlAuthenticator_LoginAs(t *testing.T) {
	t.Run("Login as valid user", func(t *testing.T) {
		authenticator, err := NewAuthenticator(validUsername, validPasswd)
		check(t, err)
		if err := authenticator.LoginAs(validUsername, validPasswd); err != nil {
			t.Error(err)
		}
	})
	t.Run("Login as invalid user", func(t *testing.T) {
		authenticator, err := NewAuthenticator(validUsername, invalidPasswd)
		check(t, err)
		if err := authenticator.LoginAs(invalidUsername, validPasswd); err == nil {
			t.Error("Expect: InvalidUsernameError\nActual: (nil)")
		}
	})
}

func TestSamlAuthenticator_SetupWith(t *testing.T) {
	t.Run("Setup with nil config", func(t *testing.T) {
		authenticator, err := NewAuthenticator(validUsername, validPasswd)
		check(t, err)
		testConf := Config{
			ShibbolethPasswordKey: DefaultPasswdKey,
			ShibbolethUsernameKey: DefaultUnameKey,
			ShibbolethAuthDomain:  DefaultAuthDomain,
			ShibbolethLoginURL:    ShibbolethLoginURL,
		}
		err = authenticator.SetupWith(testConf)
		switch e := err.(type) {
		case *ConfigDoesNotExists:
			// Expected
		default:
			t.Errorf("Expect: ConfigDoesNotExists\nActual: %+v\n", e)
		}
	})
}
