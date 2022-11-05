package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"firebase.google.com/go/v4/appcheck"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iamcredentials/v1"
)

type CustomClaims struct {
	Scope string `json:"scope"`
	Aud   string `json:"aud"`
	jwt.RegisteredClaims
}

type TokenService struct {
	serviceAccountEmail string
	tokenDuration       time.Duration
	appcheckClient      *appcheck.Client
}

func NewTokenService(serviceAccountEmail string, tokenDuration time.Duration, appcheckClient *appcheck.Client) (*TokenService, error) {
	log.Printf("App check enabled? %t", (appcheckClient != nil))
	return &TokenService{serviceAccountEmail, tokenDuration, appcheckClient}, nil
}

func (ts *TokenService) VerifyAppCheck(token string) (bool, error) {
	if ts.appcheckClient == nil {
		return false, errors.New("appcheck client not configured")
	}

	_, err := ts.appcheckClient.VerifyToken(token)
	if err == nil {
		return true, nil
	}

	return false, err
}

func (ts *TokenService) GenerateToken() (string, error) {
	claims := CustomClaims{
		"https://www.googleapis.com/auth/geo-platform.routes",
		"https://routes.googleapis.com/",
		jwt.RegisteredClaims{
			Issuer:    ts.serviceAccountEmail,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ts.tokenDuration)),
			Subject:   ts.serviceAccountEmail,
		},
	}
	b, _ := json.Marshal(claims)

	iamcredentialsServce, err := iamcredentials.NewService(context.Background())
	if err != nil {
		log.Printf("Error constructing IAMCredential Service: %s", err)
		return "", err
	}

	serviceAccountName := fmt.Sprintf("projects/-/serviceAccounts/%s", ts.serviceAccountEmail)
	resp, err := iamcredentialsServce.Projects.ServiceAccounts.SignJwt(
		serviceAccountName,
		&iamcredentials.SignJwtRequest{
			Delegates: []string{serviceAccountName},
			Payload:   string(b),
		}).Do()
	if err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			if e.Code == 403 {
				log.Printf("Authorization error. %s probably needs 'Service Account Token Creator' IAM role.", ts.serviceAccountEmail)
			} else {
				log.Printf("Error using IAMCredential Service to sign JWT: %v", err)
			}
		} else {
			log.Printf("Error signing JWT: %v", err)
		}
		return "", err
	}

	return resp.SignedJwt, nil
}
