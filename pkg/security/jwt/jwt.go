package jwtutil

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// The Claims struct holds the JWT payload.
type Claims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

var (
	// The secret key used for signing and validating the tokens.
	secretKey = []byte("super-secret-key")
	// The service that issues the tokens.
	issuer = "auth-service"
)

// GenerateToken creates a new JWT with the given user details and audience.
func GenerateToken(userID string, roles []string, audience string) (string, error) {
	// Set the token's expiration time to 24 hours from now.
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create a new claims object with our custom and standard claims.
	claims := &Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{audience},
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create a new token with the signing method and claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key and return the signed string.
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("could not sign the token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken parses and validates a JWT string.
func ValidateToken(tokenString, requiredAudience string) (*Claims, error) {
	// Define the claims and options for parsing.
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key to validate the signature.
		return secretKey, nil
	}, jwt.WithAudience(requiredAudience)) // Validate the audience claim.

	// Check for parsing or validation errors.
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Assert that the token is valid.
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}
