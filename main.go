package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
)

var config oauth1.Config

var requestToken, requestSecret string

func init() {
	config = oauth1.Config{
		ConsumerKey:    os.Getenv("GARMIN_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("GARMIN_CONSUMER_SECRET"),
		CallbackURL:    os.Getenv("GARMIN_CALLBACK_URL"),
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: "https://connectapi.garmin.com/oauth-service/oauth/request_token",
			AuthorizeURL:    "https://connect.garmin.com/oauthConfirm",
			AccessTokenURL:  "https://connectapi.garmin.com/oauth-service/oauth/access_token",
		},
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handle Main")
	var body = `<html><body><a href="/login">Login</a></body></html>`
	fmt.Fprintf(w, body)
}

func handleGarminLogin(w http.ResponseWriter, r *http.Request) {
	var err error
	requestToken, requestSecret, err = config.RequestToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Request Token failed: %s", err.Error())
	}

	authorizationURL, err := config.AuthorizationURL(requestToken)
	url := authorizationURL.String()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Authorization URL failed: %s", err.Error())
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGarminCallback turns the oauth_token and oauth_verifier into a oauth token that can be used for requests
func handleGarminCallback(w http.ResponseWriter, r *http.Request) {
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Error parsing authorization callback: %s", err.Error())
	}
	accessToken, accessSecret, err := config.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Error creating Access Token/Access Secret: %s", err.Error())
	}
	output := fmt.Sprintf("Access Token: %s\nAccess Secret: %s\n", accessToken, accessSecret)
	w.Write([]byte(output))
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGarminLogin)
	http.HandleFunc("/callback", handleGarminCallback)
	http.HandleFunc("/push/dailies", HandleDailies)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}
