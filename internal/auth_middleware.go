package internal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type UserClaims struct {
	ID     string
	Traits map[string]any
}

func withUserData(ctx context.Context, v *UserClaims) context.Context {
	return context.WithValue(ctx, "req.session", v)
}

func GetCaller(ctx context.Context) *UserClaims {
	return ctx.Value("req.session").(*UserClaims)
}

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			claims, err := parseJWT(token)
			if err != nil {
				log.Printf("rejected token %s: %v", token, err)
				// TODO: replace with just 401
				http.Redirect(w, r, "/.ory/self-service/login/browser", http.StatusSeeOther)
				return
			}

			ctx := withUserData(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type jwtClaims struct {
	UserID  string `json:"sub"`
	Session struct {
		Identity struct {
			Traits map[string]any `json:"traits"`
		} `json:"identity"`
	} `json:"session"`
}

// parseJWT deliberately DOES NOT verify the jwt token.
// this microservice is meant to be used behind an auth gateway (like ory proxy)
// so all jwts are assumed to be verified already.
// JWT is decoded from base64 and parsed here manually to highlight that absolutely no signature verification is done.
func parseJWT(token string) (*UserClaims, error) {
	parts := strings.SplitN(token, ".", 3)
	if len(parts) != 3 {
		return nil, errors.New("unable to split jwt into three parts")
	}
	byteClaims, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("unable to decode payload part from base64: %w", err)
	}

	var claims jwtClaims
	if err := json.Unmarshal(byteClaims, &claims); err != nil {
		return nil, fmt.Errorf("unable to unmarshal claims: %w", err)
	}

	return &UserClaims{
		ID:     claims.UserID,
		Traits: claims.Session.Identity.Traits,
	}, nil
}
