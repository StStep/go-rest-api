package auth

import (
	"testing"
)

func TestVerifyToken(t *testing.T) {
	tok, err := MakeToken()
        if err != nil {
                panic(err)
        }

	err = VerifyToken(tok)
        if err != nil {
		t.Errorf("Token is invalid with error: %v", err)
        }
}
