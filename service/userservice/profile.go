package userservice

import (
	"gameApp/param"
	"gameApp/pkg/richerror"
)

//all request inputs for interactor/service should be sanitized.

func (s Service) Profile(req param.ProfileRequest) (param.ProfileResponse,error) {
	const  op = "userservice.Profile"
	//getUserByID
	user, err := s.repo.GetUserByID(req.UserID)
	if err!=nil{
		return param.ProfileResponse{},richerror.New(op).WithErr(err).
		WithMeta(map[string]interface{}{"req": req})
		
	}

	return param.ProfileResponse{Name: user.Name},nil
}

