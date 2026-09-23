package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	passwords := []string{"pa$$word", "password", "dingleberries", ""}

	for _, password := range passwords {
		hash, err := HashPassword(password)
		if err != nil {
			t.Errorf("HashPassword gone done there errored: %v", err)
		}
		if hash == "" {
			t.Errorf("HashPassword gone done disappeared")
		}
	}
}

func TestCheckPasswordHash(t *testing.T) {
	type testCase struct {
		input    string
		hash     string
		expected bool
	}
	hash1, _ := HashPassword("pa$$word")
	hash2, _ := HashPassword("password")

	tests := []testCase{
		{input: "pa$$word", hash: hash1, expected: true},
		{input: "passwrod", hash: hash2, expected: false},
	}

	for _, tt := range tests {
		result, _ := CheckPasswordHash(tt.input, tt.hash)
		if result != tt.expected {
			t.Errorf("Password: %s does not match expected hash: %v ... %v", tt.input, result, tt.expected)
		}
		fmt.Printf("Expected: %v\nActual: %v\n", tt.expected, result)
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	secrets := []string{"imASecret", "i'mAlsoASecret", "imNotSoSecret", ""}

	for _, secret := range secrets {
		userID := uuid.New()
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("Tokens done broke idiot: %v\n", err)
		}
		fmt.Printf("Token created: %v\n", token)
		validatedID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("Stupid man couldn't validate me: %v\n", err)
		}
		fmt.Printf("ID validated: %v\n", validatedID)
		if userID != validatedID {
			t.Errorf("The fuck you mean userID (%s) doesn't match validatedID (%s)\n", userID, validatedID)
		}
		fmt.Printf("IDs %v & %v match, good job!\n", userID, validatedID)
	}
}

func TestInvalidJWT(t *testing.T) {
	secrets := []string{"imASecret", "i'mAlsoASecret", "imNotSoSecret", ""}

	for _, secret := range secrets {
		userID := uuid.New()
		token, err := MakeJWT(userID, secret, -time.Second)
		if err != nil {
			t.Errorf("Tokens done broke idiot: %v\n", err)
		}

		validatedID, err := ValidateJWT(token, secret)
		if err == nil {
			t.Errorf("Stupid man validated an expired token: %v\n", err)
		}

		if userID == validatedID {
			t.Errorf("Expired token somehow returned the original user ID: %s\n", userID)
		}
	}
}
