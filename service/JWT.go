package service

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	rsa *rsa.PrivateKey
}

func NewJWTService(rsaStr string) (JWT, error) {
	privKey, err := parseRSA(rsaStr)
	if err != nil {
		return JWT{}, err
	}
	return JWT{rsa: privKey}, nil
}

func parseRSA(rsaStr string) (*rsa.PrivateKey, error) {
	privPem, _ := pem.Decode([]byte(rsaStr))
	var parsedKey any
	var err error
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPem.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPem.Bytes); err != nil {
			return nil, err
		}
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not a valid rsa private key")
	}

	return privateKey, nil
}

// generates a signed JWT with the RSA algorithm
func (j *JWT) GenerateJWT(info map[string]any, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(expiry).Unix(),
	}

	for k, v := range info {
		if _, exists := claims[k]; exists {
			continue
		}
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.rsa)
}

func (j *JWT) PublicKey() any {
	fmt.Println("public key accessed")
	return j.rsa.Public()
}

// verifies JWT using RSA and returns claims
func (j *JWT) ClaimsFromJWT(token string) (map[string]any, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}

		return j.rsa.PublicKey, nil
	})

	if err != nil {
		return map[string]any{}, err
	}

	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		return claims, nil
	}

	return map[string]any{}, errors.New("Invalid JWT")
}
