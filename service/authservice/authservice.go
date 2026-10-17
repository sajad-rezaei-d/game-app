package authservice

import (
	"errors"
	"gameApp/entity"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	signKey string
	accessExpirationTime time.Duration
	refreshExpirationTime time.Duration
	accessSubject string
	refreshSubject string
}

func New(signKey, accessSubject, refreshSubject string ,
	 accessExpirationTime,refreshExpirationTime time.Duration) Service {
		return Service{
			signKey: signKey,
			accessExpirationTime: accessExpirationTime,
			refreshExpirationTime: refreshExpirationTime,
			accessSubject: accessSubject,
			refreshSubject: refreshSubject,
		}
}

func (s Service) CreateAccessToken(user entity.User) (string,error) {
	return s.createToken(user.ID ,s.accessSubject, s.accessExpirationTime)
}

func (s Service) CreateRefreshToken(user entity.User) (string,error) {
	return s.createToken(user.ID, s.refreshSubject , s.refreshExpirationTime)
}

func (s Service) ParseToken(bearerToken string) (*Claims,error) {
	
	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (any, error) {
			return []byte(s.signKey), nil
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
}


func (s Service) createToken(userID uint, subject string, expireDuration time.Duration) (string, error) {
	//TODO: use method RS256 
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		},
		UserID: userID,
	})

	return t.SignedString([]byte(s.signKey))
}