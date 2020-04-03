package main

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

const (
	CallBackURL = "http://localhost:1333/auth/callback"

	UserInfoAPIEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
	ScopeEmail          = "https://www.googleapis.com/auth/userinfo.email"
	ScopeProfile        = "https://www.googleapis.com/auth/userinfo.profile"
)

var OAuthConf *oauth2.Config

func init() {
	OAuthConf = &oauth2.Config{
		ClientID:     "826357128682-7j3osck2fn7b9o1tl3bvk384shouvrsd.apps.googleusercontent.com",
		ClientSecret: "Zd1sA9aPl0AFDrTW6zzH7RiL",
		RedirectURL:  CallBackURL,
		Scopes:       []string{ScopeEmail, ScopeProfile},
		Endpoint:     google.Endpoint,
	}
}

func GetLoginURL(state string) string {
	return OAuthConf.AuthCodeURL(state)
}

func RandToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
