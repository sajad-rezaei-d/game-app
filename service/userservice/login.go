package userservice

import (
	"fmt"
	"gameApp/param"
	"gameApp/pkg/richerror"

	"golang.org/x/crypto/bcrypt"
)


func (s Service) Login(req param.LoginRequest) (param.LoginResponse,error){
	const op= "userservice.login"
	// TODO: it would be better to user two separate method for existence check and getUserByPhoneNumber
	user,err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err!=nil{
		return param.LoginResponse{},richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	
	if bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(req.Password)) !=nil {
		return param.LoginResponse{},fmt.Errorf("username or password isn't correct")
	}

	// The first way (generate sessionID):
	// generate random session id
	// save sessionID in db for the given UserID
	// return session id to user

	// The second way (generate jwt token)
	accessToken, err := s.auth.CreateAccessToken(user)
	if err!=nil{
		return param.LoginResponse{},fmt.Errorf("unexpected error: %w",err)
	}

	refreshToken ,err := s.auth.CreateRefreshToken(user)
	if err!=nil{
		return param.LoginResponse{},fmt.Errorf("unexpected error: %w",err)
	}

	return param.LoginResponse{
		User : param.UserInfo{
			ID: user.ID,
			PhoneNumber: user.PhoneNumber,
			Name: user.Name,
		},
		Token: param.Tokens{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
		},
	},nil
}