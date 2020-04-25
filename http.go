package ghost

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

const (
	tokenAudience = "/v3/admin/"
	tokenType     = "Ghost"
)

type adminTokenSource struct {
	Key string
}

func (ats *adminTokenSource) Token() (*oauth2.Token, error) {
	split := strings.Split(ats.Key, ":")
	if len(split) != 2 {
		return nil, fmt.Errorf("incorrect key format")
	}
	kid, secret := split[0], split[1]

	secretBytes, err := hex.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("secret portion of key not valid hex")
	}

	claims := &jwt.StandardClaims{
		Audience:  tokenAudience,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Unix() + (5 * 60),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = kid
	ss, err := token.SignedString(secretBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth token: %v", err)
	}

	return &oauth2.Token{
		AccessToken: ss,
		Expiry:      time.Now().Add(time.Minute * 5),
		TokenType:   tokenType,
	}, nil
}

// NewTokenAuthClient creates a http.Client that handles JWT auth with the server.
// This creates the http client you will want for the NewAdminClient constructor.
func NewTokenAuthClient(key string) (*http.Client, error) {
	matched, _ := regexp.MatchString("[0-9a-f]{26}", key)
	if !matched {
		return nil, fmt.Errorf("key must contain 26 hexadecimal characters")
	}

	if strings.Count(key, ":") != 1 {
		return nil, fmt.Errorf("key must be split between id and key, seperated by ':'")
	}

	ts := oauth2.ReuseTokenSource(nil, &adminTokenSource{Key: key})
	httpClient := oauth2.NewClient(context.Background(), ts)
	httpClient.Timeout = time.Second * 10
	return httpClient, nil
}
