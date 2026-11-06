package backofficeuserhandler

import (
	"gameApp/service/authorizationservice"
	"gameApp/service/authservice"
	"gameApp/service/backofficeuserservice"
)

type Handler struct {
	authConfig 	  authservice.Config
	authSvc       authservice.Service
	backofficeUserSvc       backofficeuserservice.Service
	authorizationSvc	authorizationservice.Service
	
}

func New(authConfig authservice.Config,authSvc authservice.Service,
	backofficeUserSvc    backofficeuserservice.Service,
	authorizatonSvc authorizationservice.Service ) Handler {
	return Handler{
		authConfig:authConfig,
		authSvc: authSvc,
		backofficeUserSvc: backofficeUserSvc,
		authorizationSvc: authorizatonSvc,
	}
}

