package auth

import (
	"time"
	"net/http"
	"errors"
	"github.com/o1egl/paseto"
)

var symmetricKey = []byte("YELLOW SUBMARINE, BLACK WIZARDRY")

func MakeToken() (string, error) {
	now := time.Now()
	exp := now.Add(8 * time.Hour)
	nbt := now

	jsonToken := paseto.JSONToken{
		Audience:   "test",
		Issuer:     "test_service",
		Jti:        "123",
		Subject:    "test_subject",
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}
	// Add custom claim to the token
	jsonToken.Set("data", "this is a signed message")
	footer := "some footer"

	// Encrypt data
	return paseto.NewV2().Encrypt(symmetricKey, jsonToken, paseto.WithFooter(footer))
}

func VerifyToken(token string) error {
	var jsonToken paseto.JSONToken
	var footer string
	err := paseto.NewV2().Decrypt(token, symmetricKey, &jsonToken, &footer)
	if err != nil {
		return err
	}

	if time.Now().After(jsonToken.Expiration) {
		err = errors.New("Token has expired")
	}
	if jsonToken.Issuer != "test_service" {
		err = errors.New("Unknown service " + jsonToken.Issuer)
	}

	return err
}

func AuthRequest(r *http.Request) bool {
	u, p, b := r.BasicAuth()
	if !b || u != "token" {
		return false;
	}
	return VerifyToken(p) == nil
}
