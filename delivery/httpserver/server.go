package httpserver

import (
	"fmt"
	"gameApp/config"
	"gameApp/delivery/httpserver/backofficeuserhandler"
	"gameApp/delivery/httpserver/userhandler"
	"gameApp/service/authorizationservice"
	"gameApp/service/authservice"
	"gameApp/service/backofficeuserservice"
	"gameApp/service/userservice"
	"gameApp/validator/uservalidator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config        config.Config
	userHandler   userhandler.Handler
	backofficeUserHandler backofficeuserhandler.Handler

}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service, userValidator uservalidator.Validator,
	backofficeUserSvc    backofficeuserservice.Service,
	authorizatonSvc authorizationservice.Service) Server {
	return Server{
		config:  config,
		userHandler: userhandler.New(config.Auth,authSvc,userSvc,userValidator),
		backofficeUserHandler: backofficeuserhandler.New(config.Auth, authSvc,backofficeUserSvc, authorizatonSvc ),
	}
}

func (s Server) Serve() {

	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLogger()) // log requests with log/slog
	e.Use(middleware.Recover())       // recover from panics and return an error

	// Routes
	e.GET("/health-check", s.healthCheck)
	
	s.userHandler.SetRoutes(e)
	s.backofficeUserHandler.SetRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))

}
