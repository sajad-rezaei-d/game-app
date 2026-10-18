package config

import (
	"gameApp/repository/mysql"
	"gameApp/service/authservice"
)

type HTTPServer struct {
	Port int
}

type Config struct {
	HTTPServer HTTPServer
	Auth authservice.Config
	MySql mysql.Config
}

