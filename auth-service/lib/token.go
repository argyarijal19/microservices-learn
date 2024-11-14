package lib

import (
	"auth-service/model"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTCustomClaims struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	jwt.RegisteredClaims
}

var revokedTokens = make(map[string]time.Time)

func GenerateToken(userData *model.User) (string, int64, error) {
	fmt.Printf("GenerateToken struct with pointer is: %v", userData)

	expirationTime := time.Now().Add(120 * time.Minute)
	claims := JWTCustomClaims{
		ClientID:     userData.ClientID,
		ClientSecret: userData.ClientSecret,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY_JWT")))

	expiresIn := int64(time.Until(expirationTime).Seconds())
	return tokenString, expiresIn, err
}

func VerifyJWTToken(tokenString string, secretKey string) (interface{}, error) {

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	var key = []byte(secretKey)

	// Cek apakah token sudah di-revoke
	if expiration, exists := revokedTokens[tokenString]; exists {
		if time.Now().Before(expiration) {
			return nil, errors.New("token has been revoked")
		}
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid access")
		}
		return nil, fmt.Errorf("access expired: %v", err)
	}

	claims, ok := token.Claims.(*JWTCustomClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("access invalid or expired")
	}

	return claims, nil
}

func RevokeToken(tokenString string) error {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse token to check if it's valid and get the expiration time
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &JWTCustomClaims{})
	if err != nil {
		return fmt.Errorf("failed to parse token: %v", err)
	}

	// Extract claims to get the expiration time
	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok {
		return errors.New("failed to parse claims")
	}

	// Store the token in the revoked list with its expiration time
	revokedTokens[tokenString] = claims.ExpiresAt.Time
	log.Printf("Token revoked: %s, expires at %v", tokenString, claims.ExpiresAt.Time)
	return nil
}
