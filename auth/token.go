package auth

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func GenerateToken(clientUuid uuid.UUID, key *ecdsa.PrivateKey) (string, error) {
	curTime := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims{
		"iss": "squad-session-server",
		"sub": clientUuid.String(),
		"aud": "squad-session-server",
		"exp": curTime.Add(time.Hour * 24 * 30).Unix(),
		"iat": curTime.Unix(),
	})

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseAndVerifyTokenClaims(tokenString string, key *ecdsa.PublicKey) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token provided")
	}

	return claims, nil
}
