package middleware

import (
	"errors"
	"gameApp/service/authservice"
	cfg "gameApp/config"
	"github.com/golang-jwt/jwt/v5"
	mw "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func Auth(service authservice.Service,config authservice.Config) echo.MiddlewareFunc {
	return mw.WithConfig(mw.Config{
		ContextKey:cfg.AuthMiddlewareContextKey ,
		SigningKey: []byte(config.SignKey),
		//TODO: as sign method string to config...
		SigningMethod: "HS256",
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claims := &authservice.Claims{}

			token, err := jwt.ParseWithClaims(
				auth,
				claims,
				func(t *jwt.Token) (any, error) {
					return []byte(config.SignKey), nil
				},
				jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
			)
			if err != nil {
				return nil, err // شامل امضای نامعتبر و توکن منقضی‌شده
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}

			return claims, nil
		},
	})
}
