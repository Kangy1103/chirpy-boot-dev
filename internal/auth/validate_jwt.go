package auth

import (
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Printf("Invalid token: %s", err)
		return uuid.Nil, err
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("UUID not found: %s", err)
	}
	subjectUUID, err := uuid.Parse(subject)
	if err != nil {
		log.Printf("UUID not found: %s", err)
		return uuid.Nil, err
	}
	return subjectUUID, nil
}
