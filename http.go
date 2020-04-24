package ghost

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

type adminTokenSource struct {
	Key string
}

func (ats *adminTokenSource) Token() (*oauth2.Token, error) {
	split := strings.Split(ats.Key, ":")
	id, secret := split[0], split[1]

	claims := &jwt.StandardClaims{
		Id:        id,
		Audience:  "/v3/admin",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Unix() + (5 * 60),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth token: %v", err)
	}

	return &oauth2.Token{
		AccessToken: ss,
		Expiry:      time.Now().Add(time.Minute * 5),
	}, nil
}

// NewTokenHTTPClient creates an admin http client that handles JWT auth with the server
func NewTokenHTTPClient(key string) (*http.Client, error) {
	matched, _ := regexp.MatchString("[0-9a-f]{26}", key)
	if !matched {
		return nil, fmt.Errorf("key must contain 26 hexadecimal characters")
	}

	if strings.Count(key, ":") != 1 {
		return nil, fmt.Errorf("key must be split between id and key, seperated by ':'")
	}

	ts := oauth2.ReuseTokenSource(nil, &adminTokenSource{})
	httpClient := oauth2.NewClient(context.Background(), ts)
	httpClient.Timeout = time.Second * 10
	return httpClient, nil
}
