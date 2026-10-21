package main

import (
	"fmt"
	"gameApp/config"
	"gameApp/delivery/httpserver"
	//"gameApp/repository/migrator"
	"gameApp/repository/mysql"
	"gameApp/service/authservice"
	"gameApp/service/userservice"
	"time"
)

const (
	JwtSignKey = "jwt_secret"
	AccessTokenSubject ="ac"
	RefreshTokenSubject ="rt"
	AccessTokenExpireDuration = time.Hour*24 
	RefreshTokenExpireDuration = time.Hour*24*7
)

func main(){
	fmt.Println("start echo server")
	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8080},
		Auth: authservice.Config{
			SignKey: JwtSignKey,
			AccessExpirationTime: AccessTokenExpireDuration,
			RefreshExpirationTime: RefreshTokenExpireDuration,
			AccessSubject: AccessTokenSubject,
			RefreshSubject: RefreshTokenSubject,
		},
		MySql: mysql.Config{
			UserName: "gameapp",
			Password: "gameappt0lk2o20",
			Port: 3308,
			Host: "localhost",
			DBName: "gameapp-db",
		},
	}
	
	//TODO: add command for migrations
	//mgr := migrator.New(cfg.MySql)
	//mgr.Up()

	authSvc , userSvc :=setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc)	

	fmt.Println("start echo server")
	server.Serve()

}


func setupServices(cfg config.Config)(authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)
	
	MYSqlRepo := mysql.New(cfg.MySql)
	userSvc:= userservice.New(authSvc, MYSqlRepo)

	return authSvc,userSvc

}