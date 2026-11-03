package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	UserName string `koanf:"username"`
	Password string `koanf:"password"`
	Port     int    `koanf:"port"`
	Host     string `koanf:"host"`
	DBName   string `koanf:"db_name"`
}
type MYSQLDB struct {
	config Config
	db     *sql.DB
}

func New(config Config) *MYSQLDB {
	// parseTime=true changes the output type of DATE and DATETIME values to time.Time
	// instead of []byte / string
	// The date or datetime like 0000-00-00 00:00:00 is convert into zero value of time.Time
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=True",
		config.UserName, config.Password, config.Host, config.Port, config.DBName))
	if err != nil {
		panic(fmt.Errorf("can't open mysql db: %v", err))
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MYSQLDB{config: config, db: db}
}
