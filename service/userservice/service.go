package userservice

import (
	"gameApp/entity"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	IsPhoneNumberUnique (phoneNumber string) (bool,error)
	Register(u entity.User) (entity.User , error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User,bool,error)
	GetUserByID(userID uint) (entity.User,error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken(user entity.User)(string, error)
}

type Service struct{
	auth AuthGenerator
	repo Repository
}

func New(authGenerator AuthGenerator,repo Repository) Service{
	return Service{repo: repo, auth: authGenerator}
}

func getBcryptHash(password string) string{
	pass := []byte(password)
	hPass,_ :=bcrypt.GenerateFromPassword(pass, 0)
	return string(hPass)
}


