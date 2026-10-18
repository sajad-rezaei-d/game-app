package main

import (
	"fmt"
	"gameApp/config"
	"gameApp/delivery/httpserver"
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
	authSvc , userSvc :=setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc)	

	server.Serve()

}
/*
func userLoginHandler(writer http.ResponseWriter, req *http.Request){
	if req.Method!=http.MethodPost {
		fmt.Fprintf(writer,`{"error": "invalid method"}`)
	}
	
	data,err := io.ReadAll(req.Body)
	if err!=nil{
		writer.Write(fmt.Appendf(nil, `{"error": "%s"}`, err.Error()))
	}

	var lReq userservice.LoginRequest
	err = json.Unmarshal(data,&lReq)
	if err!=nil{
		writer.Write(fmt.Appendf(nil, `{"error": "%s"}`, err.Error()))

		return
	}
	authSvc := authservice.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject, AccessTokenExpireDuration, RefreshTokenExpireDuration)
	mysqlRepo :=mysql.New()
	userSvc := userservice.New(authSvc,mysqlRepo)

	resp,err := userSvc.Login(lReq)
	if err!=nil{
		writer.Write(fmt.Appendf(nil, `{"error": "%s"}`, err.Error()))

		return
	}
	
	data ,err = json.Marshal(resp)
	if err!=nil{
		writer.Write(fmt.Appendf(nil, `{"error": "%s"}`, err.Error()))

		return
	}

	
	writer.Write(data)

}


func userProfileHandler(writer http.ResponseWriter,req *http.Request) {
	if req.Method!=http.MethodGet {
		fmt.Fprintf(writer,`{"error": "invalid method"}`)
	}
	// The first way (generate sessionID) : 
	// validate sessionID by database and get userID
	// sessionID := req.Header.Get("SessionID")

	// The second way (generate jwt token) :
	authSvc := authservice.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject, AccessTokenExpireDuration, RefreshTokenExpireDuration)
	authToken:= req.Header.Get("Authorization")
	claims , err := authSvc.ParseToken(authToken)
	if err!=nil {
		fmt.Fprintf(writer, `{"error": "token is not valid"}`)

	}

	mysqlRepo :=mysql.New()
	userSvc := userservice.New(authSvc,mysqlRepo)

	
	resp,err := userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
	if err!=nil{
		writer.Write(fmt.Appendf(nil, `{"error": "%s"}`, err.Error()))

		return
	}

	data ,err := json.Marshal(resp)
	if err!=nil{
		writer.Write(fmt.Appendf(nil, `{"error": "%s"}`, err.Error()))

		return
	}

	
	writer.Write(data)
}
*/

func setupServices(cfg config.Config)(authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)
	
	MYSqlRepo := mysql.New(cfg.MySql)
	userSvc:= userservice.New(authSvc, MYSqlRepo)

	return authSvc,userSvc

}