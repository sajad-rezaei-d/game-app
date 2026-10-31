package userhandler

import (
	"gameApp/service/authservice"
	"gameApp/service/userservice"
	"gameApp/validator/uservalidator"

)

type Handler struct {
	authSvc       authservice.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
}

func New(authSvc authservice.Service, userSvc userservice.Service, userValidator uservalidator.Validator) Handler {
	return Handler{
		authSvc: authSvc,
		userSvc: userSvc,
		userValidator: userValidator,
	}
}

