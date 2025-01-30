package auth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
	Admin        TokenType = "admin"
	ApiKey       TokenType = "apiKey"
)

type IAuth interface {
	SignToken() string
}

type IAdminAuth interface {
	SignToken() string
}

type auth struct {
	mapClaims *authMapClaims
	cfg       config.IJwtConfig
}

type adminAuth struct {
	*auth
}

type authMapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func NewAuth(tokenType TokenType, cfg config.IJwtConfig, clams *users.UserClaims) (IAuth, error) {
	switch tokenType {
	case AccessToken:
		return newAccessToken(cfg, clams), nil
	case RefreshToken:
		return newRefreshToken(cfg, clams), nil
	case Admin:
		return newAdminToken(cfg), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
}

func (a *auth) SignToken() string {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		a.mapClaims,
	)
	tokenSS, _ := token.SignedString(a.cfg.SecretKey())

	return tokenSS
}

func (a *adminAuth) SignToken() string {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		a.mapClaims,
	)
	tokenSS, _ := token.SignedString(a.cfg.AdminKey())

	return tokenSS
}
func ParseAdmin(cfg config.IJwtConfig, tokenString string) (*authMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.AdminKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expired")
		} else {
			return nil, fmt.Errorf("parse token error: %w", err)
		}
	}

	if claims, ok := token.Claims.(*authMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}

}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*authMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecretKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expired")
		} else {
			return nil, fmt.Errorf("parse token error: %w", err)
		}
	}

	if claims, ok := token.Claims.(*authMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}

}

func RepeatToke(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &auth{
		cfg: cfg,
		mapClaims: &authMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "go-rest-api",
				Subject:   "refresh-token",
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}

	return obj.SignToken()

}

func newAccessToken(cfg config.IJwtConfig, clams *users.UserClaims) IAuth {
	return &auth{
		cfg: cfg,
		mapClaims: &authMapClaims{
			Claims: clams,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "go-rest-api",
				Subject:   "access-token",
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, clams *users.UserClaims) IAuth {
	return &auth{
		cfg: cfg,
		mapClaims: &authMapClaims{
			Claims: clams,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "go-rest-api",
				Subject:   "refresh-token",
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newAdminToken(cfg config.IJwtConfig) IAuth {
	return &adminAuth{
		&auth{
			cfg: cfg,
			mapClaims: &authMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "go-rest-api",
					Subject:   "admin-token",
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
