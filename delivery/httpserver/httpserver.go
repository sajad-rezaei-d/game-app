package httpserver

import (
	"fmt"
	"gameApp/config"
	"gameApp/service/authservice"
	"gameApp/service/userservice"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config config.Config
	authSvc authservice.Service
	userSvc userservice.Service
}

func New(config config.Config, authSvc authservice.Service ,userSvc userservice.Service) Server {
	return Server{config: config, authSvc: authSvc, userSvc: userSvc}
}

func (s Server) Serve() {
	
	e :=echo.New()

	// Middleware
	e.Use(middleware.RequestLogger()) // log requests with log/slog
	e.Use(middleware.Recover())       // recover from panics and return an error

	// Routes
	e.GET("/health-check", s.healthCheck )

	e.POST("/users/register", s.userRegister)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))

}