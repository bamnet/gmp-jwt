package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/appcheck"
	"github.com/jnovack/flag"
	"google.golang.org/api/oauth2/v2"

	"github.com/bamnet/gmp-jwt/apis"
	ts "github.com/bamnet/gmp-jwt/tokens"
)

const appCheckTokenHeader = "X-Firebase-AppCheck"
const defaultTokenDuration = 30 * time.Minute

var appCheckEnabled bool
var corsAllowedOrigin string
var tokenService *ts.TokenService
var allowedAPIs []string

// handler makes jwts.
func handler(w http.ResponseWriter, r *http.Request) {
	// Parse the requested APIs, if any.
	r.ParseForm()
	requestedAPIs := strings.Split(r.Form.Get("apis"), ",")
	if len(requestedAPIs) == 0 || requestedAPIs[0] == "" {
		requestedAPIs = allowedAPIs
	}
	// Don't let users request an un-allowed API.
	if !apis.IsWildcard(allowedAPIs) {
		requestedAPIs = filterAllowedAPIs(requestedAPIs)
	}

	if corsAllowedOrigin != "" {
		w.Header().Set("Access-Control-Allow-Origin", corsAllowedOrigin)
	}
	w.Header().Set("Access-Control-Allow-Headers", appCheckTokenHeader)
	if appCheckEnabled {
		if ok, err := tokenService.VerifyAppCheck(r.Header.Get(appCheckTokenHeader)); !ok || err != nil {
			log.Printf("App check token verification failed (%t): %v", ok, err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	token, err := tokenService.GenerateToken(requestedAPIs)
	if err != nil {
		switch err {
		case apis.ErrIncompatibleAPIs:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Write([]byte(token))
}

// whoami identified the service account email address of the available service account.
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
	flag.BoolVar(&appCheckEnabled, "enable_appcheck", false,
		fmt.Sprintf("If set, requests must have a valid token from app check in the %s header", appCheckTokenHeader))

	flag.StringVar(&corsAllowedOrigin, "cors_origins", "",
		"Value to set for the 'Access-Control-Allow-Origin' header")

	var tokenDuration time.Duration
	flag.DurationVar(&tokenDuration, "token_duration", defaultTokenDuration,
		"Duration a generated token is valid for")

	allowedAPIsFlat := flag.String("allowed_apis", "*", "comma-seperated list of APIs tokens can be generated for, or * for all supported")

	flag.Parse()

	// Parse and validate the allowed APIs flag.
	allowedAPIs = strings.Split(*allowedAPIsFlat, ",")
	if !apis.IsWildcard(allowedAPIs) {
		for _, api := range allowedAPIs {
			if _, ok := apis.APIs[api]; !ok {
				log.Printf("Warning: API '%s' not found in API Info table. Will be skipped.", api)
			}
		}
	}

	var appcheckClient *appcheck.Client
	if appCheckEnabled {
		fb, err := firebase.NewApp(context.Background(), nil)
		if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
		}

		appcheckClient, err = fb.AppCheck(context.Background())
		if err != nil {
			log.Fatalf("Error initializing AppCheck client: %v", err)
		}
	}

	email, err := whoami(context.Background())
	if err != nil {
		log.Fatalf("Error identifying service account: %v", err)
	}

	tokenService, err = ts.NewTokenService(email, tokenDuration, appcheckClient)
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func filterAllowedAPIs(requested []string) []string {
	safe := []string{}
	for _, api := range requested {
		if contains(allowedAPIs, api) {
			safe = append(safe, api)
		} else {
			log.Printf("API '%s' not in allowlist, skipping", api)
		}
	}
	return safe
}
