package config

import (
	"gameApp/repository/mysql"
	"gameApp/service/authservice"
)

type HTTPServer struct {
	Port int  `koanf:"port"`
}

type Config struct {
	HTTPServer HTTPServer	`koanf:"http_server"`
	Auth authservice.Config	`koanf:"auth"`
	MySql mysql.Config		`koanf:"mysql"`
}

