package main

import (
	"encoding/json"
	"fmt"
	"gameApp/repository/mysql"
	"gameApp/service/userservice"
	"io"
	"log"
	"net/http"
)

const (
	JwtSignKey = "jwt_secret"
)

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/health-check", healthCheckHandler)
	mux.HandleFunc("/users/register", userRegisterHandler)
	mux.HandleFunc("/users/login", userLoginHandler)
	mux.HandleFunc("/users/profile", userProfileHandler)

	server:= http.Server{ Addr: ":8080", Handler: mux}
	log.Fatal(server.ListenAndServe())

}

func userRegisterHandler(writer http.ResponseWriter, req *http.Request){
	if req.Method!=http.MethodPost {
		fmt.Fprintf(writer,`{"error": "invalid method"}`)
	}

	data,err := io.ReadAll(req.Body)
	if err!=nil{
		writer.Write(fmt.Appendf(nil, `{"error": "%s"}`, err.Error()))
	}

	var uReq userservice.RegisterRequest
	err = json.Unmarshal(data,&uReq)
	if err!=nil{
		writer.Write(fmt.Appendf(nil, `{"error": "%s"}`, err.Error()))

		return
	}

	mysqlRepo :=mysql.New()
	userSvc := userservice.New(mysqlRepo, JwtSignKey)

	_,err = userSvc.Register(uReq)
	if err!=nil{
		writer.Write(fmt.Appendf(nil, `{"error": "%s"}`, err.Error()))

		return
	}
	
	writer.Write([]byte(`{"message": "user created"}`))
}

func healthCheckHandler(writer http.ResponseWriter,req *http.Request){
	fmt.Fprintf(writer,`{"message": "everyting is good!"}`)
}


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

	mysqlRepo :=mysql.New()
	userSvc := userservice.New(mysqlRepo, JwtSignKey)

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

	

	var pReq userservice.ProfileRequest
	
	
	data,err := io.ReadAll(req.Body)
	if err!=nil{
		writer.Write(fmt.Appendf(nil, `{"error": "%s"}`, err.Error()))
	}

	err = json.Unmarshal(data,&pReq)
	if err!=nil{
		writer.Write(fmt.Appendf(nil, `{"error": "%s"}`, err.Error()))

		return
	}

	mysqlRepo :=mysql.New()
	userSvc := userservice.New(mysqlRepo,JwtSignKey)

	
	resp,err := userSvc.Profile(pReq)
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