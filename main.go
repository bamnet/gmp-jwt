package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/oauth2/v2"
)

type CustomClaims struct {
	Scope string `json:"scope"`
	Aud   string `json:"aud"`
	jwt.RegisteredClaims
}

var serviceAccountEmail string

func handler(w http.ResponseWriter, _ *http.Request) {
	claims := CustomClaims{
		"https://www.googleapis.com/auth/geo-platform.routes",
		"https://routes.googleapis.com/",
		jwt.RegisteredClaims{
			Issuer:    serviceAccountEmail,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Subject:   serviceAccountEmail,
		},
	}
	b, _ := json.Marshal(claims)

	iamcredentialsServce, err := iamcredentials.NewService(context.Background())
	if err != nil {
		log.Fatalf("Error constructing iam service: %s", err)
	}

	serviceAccountName := fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccountEmail)
	resp, err := iamcredentialsServce.Projects.ServiceAccounts.SignJwt(
		serviceAccountName,
		&iamcredentials.SignJwtRequest{
			Delegates: []string{serviceAccountName},
			Payload:   string(b),
		}).Do()
	if err != nil {
		log.Fatalf("Error calling signJwt: %+v", err)
	}

	w.Write([]byte(resp.SignedJwt))
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
	serviceAccountEmail = email

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
