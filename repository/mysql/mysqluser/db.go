package mysqluser

import "gameApp/repository/mysql"



type DB struct {
	conn *mysql.MYSQLDB
}

func New(conn *mysql.MYSQLDB) *DB {
	return &DB{
		conn: conn ,
	}
}