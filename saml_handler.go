package kitwalk

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/url"
)

func isContinueRequired(body io.Reader) bool {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		panic("Cannot parse response")
	}
	unameInput := doc.Find("input[id='username']").First()
	passwdInput := doc.Find("input[id='password']").First()
	if unameInput.Length() == 0 && passwdInput.Length() == 0 {
		return true
	}
	return false
}

func parseSamlResp(body io.Reader) (string, url.Values, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		panic("Cannot parse response")
	}
	// When invalid auth info is posted
	errorForm := doc.Find("p[class=\"form-error\"]").First()
	if errorForm.Length() != 0 {
		return "", nil, &ShibbolethAuthError{errMsg: errorForm.Text()}
	}
	// Parse SAML response form.
	// When you use a normal browser such as Chrome or FireFox, this form will be submitted automatically.
	form := doc.Find("form")
	actionUrl, aUrlExists := form.Attr("action")
	if !aUrlExists {
		return "", nil, &ShibbolethAuthError{errMsg: "Could not find action url"}
	}
	relayState, rStateExists := form.Find("input[name=\"" + DefaultRelayStateKey + "\"]").First().Attr("value")
	samlResponse, sRespExists := form.Find("input[name=\"" + DefaultSAMLResponseKey + "\"]").First().Attr("value")
	if !rStateExists || !sRespExists {
		return "", nil, &ShibbolethAuthError{errMsg: "Could not parse response"}
	}
	// Create post data
	authData := url.Values{}
	authData.Add(DefaultRelayStateKey, relayState)
	authData.Add(DefaultSAMLResponseKey, samlResponse)
	return actionUrl, authData, nil
}
