package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

func TestParseRSA(t *testing.T) {
	t.Parallel()
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Fatalf("Unable to generate rsa %s", err.Error())
	}
	rsaBytes := x509.MarshalPKCS1PrivateKey(key)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: rsaBytes,
		},
	)
	rsaString := string(privkey_pem)
	resultKey, err := parseRSA(rsaString)
	if err != nil {
		t.Errorf("Error parsing rsa key %s", err.Error())
	}
	if !resultKey.Equal(key) {
		t.Error("Parsed key does not equal original key")
	}
}

func TestClaimsFromJWT(t *testing.T) {
	t.Parallel()
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Fatalf("Unable to generate rsa %s", err.Error())
	}
	rsaBytes := x509.MarshalPKCS1PrivateKey(key)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: rsaBytes,
		},
	)
	rsaString := string(privkey_pem)
	claims := jwt.MapClaims{
		"randomClaim": "15",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwt, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("Unable to generate jwt %s", err.Error())
	}

	jwtService, err := NewJWTService(rsaString)
	resultClaims, err := jwtService.ClaimsFromJWT(jwt)

	if err != nil {
		t.Errorf("Error getting claims from jwt %s", err.Error())
	}

	for k, v := range claims {
		if resultValue, exists := resultClaims[k]; !exists || !reflect.DeepEqual(resultValue, v) {
			t.Error("Error jwt claims do not match")
		}
	}
}
