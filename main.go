package main

import (
	"fmt"
	"gameApp/config"
	"gameApp/delivery/httpserver"
	"gameApp/repository/migrator"
	"gameApp/repository/mysql"
	"gameApp/repository/mysql/mysqlaccesscontrol"
	"gameApp/repository/mysql/mysqluser"
	"gameApp/service/authorizationservice"
	"gameApp/service/authservice"
	"gameApp/service/backofficeuserservice"
	"gameApp/service/userservice"
	"gameApp/validator/uservalidator"
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
	//TODO: reaad config path from command line
	cfg2:= config.Load("config.yml")
	fmt.Printf("cfg2: %+v\n",cfg2)
	//TODO: merge cfg with cfg2	
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
			Port: 3309,
			Host: "localhost",
			DBName: "gameapp-db",
		},
	}
	
	//TODO: add command for migrations
	mgr := migrator.New(cfg.MySql)
	mgr.Up()

	authSvc , userSvc, userValidator, backofficeSvc, authorizationSvc :=setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc)	

	fmt.Println("start echo server")
	server.Serve()

}


func setupServices(cfg config.Config)(authservice.Service, userservice.Service, uservalidator.Validator,
	 backofficeuserservice.Service,authorizationservice.Service) {
	authSvc := authservice.New(cfg.Auth)
	MYSqlRepo := mysql.New(cfg.MySql)

	userMysql := mysqluser.New(MYSqlRepo)
	userSvc:= userservice.New(authSvc, userMysql)

	aclMysql := mysqlaccesscontrol.New(MYSqlRepo)
	backofficeUserSvc := backofficeuserservice.New()

	authorizationSvc := authorizationservice.New(aclMysql)

	uv:= uservalidator.New(userMysql)
	return authSvc,userSvc,uv, backofficeUserSvc, authorizationSvc

}