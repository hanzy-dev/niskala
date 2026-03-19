package authjwt

import (
	"context"
	"fmt"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type Verifier struct {
	jwks     keyfunc.Keyfunc
	issuer   string
	audience string
	enabled  bool
}

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewVerifier(ctx context.Context, jwksURL string, issuer string, audience string) (*Verifier, error) {
	if jwksURL == "" || issuer == "" || audience == "" {
		return &Verifier{
			enabled: false,
		}, nil
	}

	jwks, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("create jwks verifier: %w", err)
	}

	return &Verifier{
		jwks:     jwks,
		issuer:   issuer,
		audience: audience,
		enabled:  true,
	}, nil
}

func (v *Verifier) Enabled() bool {
	return v != nil && v.enabled
}

func (v *Verifier) ParseBearerToken(tokenString string) (*Claims, error) {
	if !v.Enabled() {
		return nil, fmt.Errorf("jwt verifier is disabled")
	}

	tokenString = strings.TrimSpace(tokenString)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, v.jwks.Keyfunc)
	if err != nil {
		return nil, fmt.Errorf("parse jwt: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid jwt")
	}

	if claims.Issuer != v.issuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	audiences, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf("get audience: %w", err)
	}

	audienceMatched := false
	for _, aud := range audiences {
		if aud == v.audience {
			audienceMatched = true
			break
		}
	}

	if !audienceMatched {
		return nil, fmt.Errorf("invalid audience")
	}

	if claims.Subject == "" {
		return nil, fmt.Errorf("missing subject")
	}

	return claims, nil
}
