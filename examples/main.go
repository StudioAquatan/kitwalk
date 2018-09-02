package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/StudioAquatan/kitwalk"
)

func main() {
	// Create http client
	c := http.DefaultClient

	// Prepare parent context with cancel method (optional)
	// Here is an example of canceling context after 10 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create authenticator with your username and password
	// NOTE: In this step, login is not actually done yet.
	auth, err := kitwalk.NewAuthenticator(ctx, "your username", "your password")
	if err != nil {
		panic(err)
	}

	// Handle login with http.Client.
	err = auth.LoginWith(c)
	if err != nil {
		// It returns an error if an error occurs during login attempt.
		panic(err)
	}

	// Now, you can reach the website protected with Shibboleth.
	resp, err := c.Get("https://portal.student.kit.ac.jp/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// 200 https://portal.student.kit.ac.jp/
	fmt.Println(resp.StatusCode, resp.Request.URL)
}
