package internal

import (
	"context"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"time"

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

type AuthMiddleware struct {
	cfg    AuthConfig
	jwkSet jwk.Set
}

func NewAuthMiddleware(cfg AuthConfig) *AuthMiddleware {
	return &AuthMiddleware{
		cfg,
		getJWK(cfg.JWKsURL),
	}
}

func (mw *AuthMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			claims, ok := mw.validateJWT(token)
			if !ok {
				// TODO: replace with just 401
				http.Redirect(w, r, mw.cfg.LoginURL, http.StatusSeeOther)
				return
			}

			ctx := withUserData(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getJWK(jwkURL string) jwk.Set {
	c := jwk.NewCache(context.Background())
	c.Register(jwkURL, jwk.WithMinRefreshInterval(15*time.Minute))
	_, err := c.Refresh(context.Background(), jwkURL)
	if err != nil {
		log.Panicf("unable to get jwks for authentication: %v", err)
	}
	return jwk.NewCachedSet(c, jwkURL)
}

func (mw *AuthMiddleware) validateJWT(token string) (*UserClaims, bool) {
	validated, err := jwt.Parse([]byte(token), jwt.WithKeySet(mw.jwkSet, jws.WithRequireKid(false)), jwt.WithValidate(true))
	if err != nil {
		return nil, false
	}
	validated.PrivateClaims()
	sub := validated.Subject()
	traits := validated.PrivateClaims()["session"].(map[string]any)["identity"].(map[string]any)["traits"]
	return &UserClaims{
		ID:     sub,
		Traits: traits.(map[string]any),
	}, true
}
