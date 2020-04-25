package ghost

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

const (
	tokenAudience = "/v3/admin/"
	tokenType     = "Ghost"
	timeout       = time.Second * 10
)

// AdminTokenSource is a token source for token-based authentication with
// the Ghost Admin API.
type AdminTokenSource struct {
	Key string
}

// Token returns the Ghost jwt token needed for token based authenication.
func (ats *AdminTokenSource) Token() (*oauth2.Token, error) {
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

// NewAdminTokenSource returns a reusable oauth2.TokenSource that is backed by
// the AdminTokenSource implementation. It handles properly creating and renewing
// the JWT needed for communication with Ghost for token-based auth.
func NewAdminTokenSource(key string) (oauth2.TokenSource, error) {
	matched, _ := regexp.MatchString("[0-9a-f]{26}", key)
	if !matched {
		return nil, fmt.Errorf("key must contain 26 hexadecimal characters")
	}
	if strings.Count(key, ":") != 1 {
		return nil, fmt.Errorf("key must be split between id and secret, seperated by ':'")
	}

	ts := oauth2.ReuseTokenSource(nil, &AdminTokenSource{Key: key})
	return ts, nil
}
