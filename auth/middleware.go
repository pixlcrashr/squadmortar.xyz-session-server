package auth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func parseAuthorizationToken(value string) (string, error) {
	parts := strings.Split(value, "Bearer ")

	if len(parts) != 2 {
		return "", errors.New("invalid authorization header value")
	}

	return parts[1], nil
}

func ForContext(ctx context.Context) (uuid.UUID, error) {
	v := ctx.Value("clientUuid")

	u, ok := v.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("client uuid not found")
	}

	return u, nil
}

func getClientUuidFromToken(token string, key *ecdsa.PublicKey) (uuid.UUID, error) {
	claims, err := ParseAndVerifyTokenClaims(token, key)
	if err != nil {
		return uuid.UUID{}, err
	}

	clientUuidStr, ok := claims["sub"].(string)
	if !ok {
		return uuid.UUID{}, errors.New("invalid sub field in token provided")
	}

	clientUuid, err := uuid.Parse(clientUuidStr)
	if err != nil {
		return uuid.UUID{}, err
	}

	return clientUuid, nil
}

func AuthenticationMiddleware(key *ecdsa.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string
			var err error

			switch r.Header.Get("Upgrade") {
			case "websocket":
				token = r.URL.Query().Get("token")
			default:
				token, err = parseAuthorizationToken(r.Header.Get("Authorization"))
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
			}

			clientUuid, err := getClientUuidFromToken(token, key)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), "clientUuid", clientUuid)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			return
		})
	}
}
