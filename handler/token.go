package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWT struct {
	privateKey []byte
	publicKey  []byte
}

func NewJWT(privateKey []byte, publicKey []byte) JWT {
	return JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (j *JWT) GenerateToken(userID int64) (string, error) {
	// Parse private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", err
	}

	// Define token claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create the token object with RS256 algorithm
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *JWT) ValidateToken(headerAuthorization string) (jwt.MapClaims, error) {
	// Check token type
	authorization := strings.Split(headerAuthorization, " ")
	if len(authorization) == 0 || len(authorization) > 2 {
		return nil, fmt.Errorf("invalid header")
	}
	if authorization[0] != "Bearer" {
		return nil, fmt.Errorf("invalid token type")
	}
	token := authorization[1]

	// Parse public key
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return nil, err
	}

	// Parse token
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, err
	}

	// Extract the token claims
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
