package utils

import (
	"crypto"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var privKey crypto.PrivateKey
var pubKey crypto.PublicKey

func JwtInit() {
	filePrivateKeyPem, err := os.ReadFile("./keys/jwt-private.pem")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open jwt private key:\n|		%s\n", err.Error())
		os.Exit(1)
	}

	filePublicKeyPem, err2 := os.ReadFile("./keys/jwt-public.pem")
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "cannot open jwt public key:\n|		%s\n", err2.Error())
		os.Exit(1)
	}

	var err3 error
	privKey, err3 = jwt.ParseEdPrivateKeyFromPEM(filePrivateKeyPem)
	if err3 != nil {
		fmt.Fprintf(os.Stderr, "cannot parse jwt private key:\n|		%s\n", err3.Error())
		os.Exit(1)
	}
	pubKey, err3 = jwt.ParseEdPublicKeyFromPEM(filePublicKeyPem)
	if err3 != nil {
		fmt.Fprintf(os.Stderr, "cannot parse jwt public key:\n|		%s\n", err3.Error())
		os.Exit(1)
	}
}

type jwtClaim struct {
	Username string `json:"username"`
	UserType string `json:"userType"`
	jwt.RegisteredClaims
}

func JwtSignToken(username string, userType string) (string, error) {
	claims := jwtClaim{
		username,
		userType,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(privKey)
}

func JwtValidateToken(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("invalid token signing")
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
