package auth

import (
    "testing"
	"github.com/google/uuid"
	"time"
)

func TestJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := time.Hour

	JWTtoken, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}
    
	if JWTtoken == "" {
        t.Fatal("expected a non-empty token string")
    }

	verifiedID, err := ValidateJWT(JWTtoken, tokenSecret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	if verifiedID != userID {
		t.Errorf("expected userID %v, got %v", userID, verifiedID)
	}
}