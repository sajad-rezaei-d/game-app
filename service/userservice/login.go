package userservice

import (
	"fmt"
	"gameApp/dto"
	"gameApp/pkg/richerror"

	"golang.org/x/crypto/bcrypt"
)


func (s Service) Login(req dto.LoginRequest) (dto.LoginResponse,error){
	const op= "userservice.login"
	// TODO: it would be better to user two separate method for existence check and getUserByPhoneNumber
	user,exist,err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err!=nil{
		return dto.LoginResponse{},richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	if !exist {
		return dto.LoginResponse{},fmt.Errorf("record not found")
	}


	if bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(req.Password)) !=nil {
		return dto.LoginResponse{},fmt.Errorf("username or password isn't correct")
	}

	// The first way (generate sessionID):
	// generate random session id
	// save sessionID in db for the given UserID
	// return session id to user

	// The second way (generate jwt token)
	accessToken, err := s.auth.CreateAccessToken(user)
	if err!=nil{
		return dto.LoginResponse{},fmt.Errorf("unexpected error: %w",err)
	}

	refreshToken ,err := s.auth.CreateRefreshToken(user)
	if err!=nil{
		return dto.LoginResponse{},fmt.Errorf("unexpected error: %w",err)
	}

	return dto.LoginResponse{
		User : dto.UserInfo{
			ID: user.ID,
			PhoneNumber: user.PhoneNumber,
			Name: user.Name,
		},
		Token: dto.Tokens{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
		},
	},nil
}