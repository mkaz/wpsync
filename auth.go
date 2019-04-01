package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/automattic/go/jaguar"
)

// runSetup prompts the user for the necessary info to
// configure and run. It can be triggered directly using
// --init or will get triggered if testSetup fails
func runSetup() {
	var user, pass string

	// prompt user for site
	fmt.Print("Enter URL for site: ")
	_, err := fmt.Scanf("%s", &conf.SiteURL)
	if err != nil {
		log.Fatal("What happened?", err)
	}

	// prompt for username
	fmt.Print("Enter username: ")
	_, err = fmt.Scanf("%s", &user)
	if err != nil {
		log.Fatal("What happened?", err)
	}

	// prompt for password
	fmt.Print("Enter password: ")
	_, err = fmt.Scanf("%s", &pass)
	if err != nil {
		log.Fatal("What happened?", err)
	}

	// make JWT call to fetch token
	url := strings.Join([]string{conf.SiteURL, "wp-json", "jwt-auth/v1/token"}, "/")
	j := jaguar.New()
	j.Url(url)
	j.Params.Add("username", user)
	j.Params.Add("password", pass)
	resp, err := j.Method("POST").Send()
	if err != nil {
		log.Warn("API error authentication", err)
	}

	if err := json.Unmarshal(resp.Bytes, &conf); err != nil {
		log.Warn("Error parsing JSON response", string(resp.Bytes), err)
	}
}

// testSetup confirms everything is configured and working
// includes the local directories, blog config, and auth
func testSetup() bool {
	return false
}
