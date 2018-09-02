# StudioAquatan/kitwalk

[![CircleCI](https://circleci.com/gh/StudioAquatan/kitwalk.svg?style=svg)](https://circleci.com/gh/StudioAquatan/kitwalk) [![codecov](https://codecov.io/gh/StudioAquatan/kitwalk/branch/master/graph/badge.svg)](https://codecov.io/gh/StudioAquatan/kitwalk) [![Go Report Card](https://goreportcard.com/badge/github.com/StudioAquatan/kitwalk)](https://goreportcard.com/report/github.com/StudioAquatan/kitwalk)

This package will support your automation in your student life with a program written in golang.

Package `http/net` is very useful and generally used in many situations. `http.Client` is a standard method to access website or send API request ...etc, so `kitwalk` is just a wrapper to get authenticated cookies and store it to your own `http.Client`.

Just create `http.Client` instance, and give it to this package, now you can access the website protected with Shibboleth using `http.Client`.

## Environment

- Go 1.10 or later

Other packages

- [github.com/PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)

## Usage

See [example](./examples/main.go)

**NOTE**

Shibboleth authentication will revoke after few hours. If the website require to authenticate again, you should use `LoginWith` again. If you logged in and the website require to re-authenticate, `LoginWith` cannot works well, I think. So, this package only support

1. **Short life cycle work**
    - For example, scraping every minutes, CLI tool, ... etc.
    - You should clear the cookie properly and periodically if your program run for a long time.
2. **No re-authentication**
    - I don't know such a situation will occur, kitwalk doesn't support that.
    - For example, in other services, re-auth is required when you change password or perform administrative activity.

Please be careful to use. **DON'T PUSH YOUR USERNAME OR PASSWORD TO YOUR REPOSITORY.**

## Development

We welcome your contribution! Please commit any fix or improvement.

To fix version of packages, you should use `dep`.

- Dep v0.5.0 or later

### Release

1. Write documents
    - examples/main.go
    - README.md
2. `git commit -a -m 'Release v${VERSION}'`
3. `git tag v${VERSION}`
4. `git push origin master --tags`

## License

GPL v3

## Author

- StudioAquatan
    - pudding
