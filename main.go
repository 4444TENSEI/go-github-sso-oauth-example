/*
 * Port of server.rb from GitHub "Basics of authentication" developer guide
 * https://developer.github.com/v3/guides/basics-of-authentication/
 *
 * Simple OAuth server retrieving all email adresses of the GitHub user who authorizes this GitHub OAuth Application
 */

// Simple OAuth server retrieving the email adresses of a GitHub user.
package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Simple OAuth server retrieving the email addresses of a GitHub user.

type Config struct {
	GH_BASIC_CLIENT_ID string
	GH_BASIC_SECRET_ID string
	Port               string
}

var config Config

type IndexPageData struct {
	ClientId string
}

type BasicPageData struct {
	User   *github.User
	Emails []*github.UserEmail
}

type Access struct {
	AccessToken string `json:"access_token"`
	Scope       string
}

var indexPage = template.Must(template.New("index.tmpl").ParseFiles("views/index.tmpl"))
var basicPage = template.Must(template.New("basic.tmpl").ParseFiles("views/basic.tmpl"))

func main() {
	// Load TOML configuration file
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatalf("Error loading config.toml file: %v", err)
	} else {
		log.Println("config.toml loaded successfully")
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/callback", basic)
	log.Printf("Server starting on port %s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	indexPageData := IndexPageData{ClientId: config.GH_BASIC_CLIENT_ID}
	if err := indexPage.Execute(w, indexPageData); err != nil {
		log.Println(err)
	}
}

func basic(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	values := url.Values{
		"client_id":     {config.GH_BASIC_CLIENT_ID},
		"client_secret": {config.GH_BASIC_SECRET_ID},
		"code":          {code},
		"accept":        {"json"},
	}

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(values.Encode()))
	if err != nil {
		log.Print(err)
		return
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Retrieving access token failed: ", resp.Status)
		return
	}

	var access Access
	if err := json.NewDecoder(resp.Body).Decode(&access); err != nil {
		log.Println("JSON-Decode-Problem: ", err)
		return
	}

	if access.Scope != "user:email" {
		log.Println("Wrong token scope: ", access.Scope)
		return
	}

	ctx := context.Background()
	client := getGitHubClient(access.AccessToken)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Println("Could not list user details: ", err)
		return
	}

	emails, _, err := client.Users.ListEmails(ctx, nil)
	if err != nil {
		log.Println("Could not list user emails: ", err)
		return
	}

	basicPageData := BasicPageData{User: user, Emails: emails}
	if err := basicPage.Execute(w, basicPageData); err != nil {
		log.Println(err)
	}
}

// Authenticates GitHub Client with provided OAuth access token
func getGitHubClient(accessToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
