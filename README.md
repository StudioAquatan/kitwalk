# StudioAquatan/kitwalk

[![CircleCI](https://circleci.com/gh/StudioAquatan/kitwalk.svg?style=svg)](https://circleci.com/gh/StudioAquatan/kitwalk)

[![codecov](https://codecov.io/gh/StudioAquatan/kitwalk/branch/master/graph/badge.svg)](https://codecov.io/gh/StudioAquatan/kitwalk)

[![Go Report Card](https://goreportcard.com/badge/github.com/StudioAquatan/kitwalk)](https://goreportcard.com/report/github.com/StudioAquatan/kitwalk)

This package will support your automation in your student life with a program written in golang.

Package `http/net` is very useful and generally used in many situations. `http.Client` is a standard method to access website or send API request ...etc, so `kitwalk` is just a wrapper to get authenticated cookies and store it to your own `http.Client`.

Just create `http.Client` instance, and give it to this package, now you can access the website protected with Shibboleth using `http.Client`.

## Environment

- Go 1.10 or later
- Dep v0.5.0 or later

Other packages

- [github.com/PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)

## Usage

```go
package main

import (
    "net/http"

    "github.com/StudioAquatan/kitwalk
)

func main() {
    // Create client
    client := http.DefaultClient
    // Create authenticator
    // NOTE: Authentication has not been executed yet in this step.
    authenticator, err := NewAuthenticator("your username", "your password")
    if err != nil {
        panic(err)
    }
    // Login
    // If you use http.DefaultClient, you just set nil to `LoginWith`.
    err = authenticator.LoginWith(client)
    if err != nil {
        // Auth fail or server is down or ... etc
        panic(err)
    }
    // Success to auth!
    resp, err := client.Get("https://portal.student.kit.ac.jp/")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    fmt.Println(resp.Request.StatusCode)
}
```


## License

GPL v3

## Author

- StudioAquatan
    - pudding
