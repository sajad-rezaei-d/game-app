package userservice

import (
	"fmt"
	"gameApp/dto"
	"gameApp/entity"
)


func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse,error) {
// TODO: we should verify phone number by verification code 

	//create new user in storage 
	user := entity.User{
		ID: 0,
		PhoneNumber: req.PhoneNumber,
		Name: req.Name,
		Password: getBcryptHash(req.Password),
	}

	createdUser,err:=s.repo.Register(user)
	if err!=nil{
		return dto.RegisterResponse{},fmt.Errorf("unexpected error: %w",err)
	}

	//return created user
	return dto.RegisterResponse{User: dto.UserInfo{
		ID: createdUser.ID,
		PhoneNumber: createdUser.PhoneNumber,
		Name: createdUser.Name,
	}},nil
}

