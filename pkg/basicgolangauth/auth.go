package basicgolangauth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/TurterDev/basicgolang/config"
	"github.com/TurterDev/basicgolang/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

type IBasicAuth interface {
	SignToken() string
}

type IBasicAdmin interface {
	SignToken() string
}

type basicAuth struct {
	mapClaims *basicMapClaims
	cfg       config.IJwtConfig
}

type basicAdmin struct {
	*basicAuth
}

type basicMapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

// คำนวณวันหมดอายุ token
func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

// แปลง type int64 jwt to time
func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func (a *basicAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func (a *basicAdmin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*basicMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &basicMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signin method is invalid")
		}
		return cfg.SecretKey(), nil
	})
	// check error
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expire")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	//check payload
	if claims, ok := token.Claims.(*basicMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*basicMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &basicMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signin method is invalid")
		}
		return cfg.AdminKey(), nil
	})
	// check error
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expire")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	//check payload
	if claims, ok := token.Claims.(*basicMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

// func RepeatToken
func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &basicAuth{
		cfg: cfg,
		mapClaims: &basicMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "basicgolang-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func NewBasicgolangAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IBasicAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	//Generate Refresh Token
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
	return nil, nil
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IBasicAuth {
	return &basicAuth{
		cfg: cfg,
		mapClaims: &basicMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "basicgolang-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IBasicAuth {
	return &basicAuth{
		cfg: cfg,
		mapClaims: &basicMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "basicgolang-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newAdminToken(cfg config.IJwtConfig) IBasicAuth {
	return &basicAdmin{
		&basicAuth{
			cfg: cfg,
			mapClaims: &basicMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "basicgolang-api",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
