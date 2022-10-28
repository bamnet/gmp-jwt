package main

import (
	"context"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/appcheck"

	"google.golang.org/api/oauth2/v2"
)

const appCheckTokenHeader = "X-Firebase-AppCheck"

var tokenService *TokenService

func handler(w http.ResponseWriter, r *http.Request) {
	if appCheckEnabled() {
		if ok, err := tokenService.VerifyAppCheck(r.Header.Get(appCheckTokenHeader)); !ok || err != nil {
			log.Printf("App check token verification failed (%t): %v", ok, err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	token, _ := tokenService.GenerateToken()
	w.Write([]byte(token))
}

func whoami(ctx context.Context) (string, error) {
	oauth2Service, err := oauth2.NewService(ctx)
	if err != nil {
		return "", err
	}
	res, err := oauth2Service.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return "", err
	}

	log.Printf("I am %s.", res.Email)
	return res.Email, nil
}

func main() {
	email, err := whoami(context.Background())
	if err != nil {
		log.Fatalf("Error identifying service account: %v", err)
	}

	var appcheckClient *appcheck.Client
	if appCheckEnabled() {
		fb, err := firebase.NewApp(context.Background(), nil)
		if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
		}

		appcheckClient, err = fb.AppCheck(context.Background())
		if err != nil {
			log.Fatalf("Error initializing AppCheck client: %v", err)
		}
	}

	tokenService, err = NewTokenService(email, appcheckClient)
	if err != nil {
		log.Fatalf("Error creating token service: %v", err)
	}

	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func appCheckEnabled() bool {
	return os.Getenv("ENABLE_APPCHECK") == "TRUE"
}
