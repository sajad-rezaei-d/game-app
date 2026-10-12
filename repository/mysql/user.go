package mysql

import (
	"database/sql"
	"fmt"
	"gameApp/entity"
)

func (d MYSQLDB) IsPhoneNumberUnique(phoneNumber string) (bool,error){
	 user := entity.User{}
	 var createdAt []uint8

	 row := d.db.QueryRow(`select * from users where phone_number = ?`,phoneNumber)
	
	 //It must be in the exact same order as the columns defined in the database.
	 //We also need to pass its pointer so that the scan function can set the values.
	 err :=row.Scan(&user.ID,&user.Name,&user.PhoneNumber,&createdAt)
	 if err!=nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false , fmt.Errorf("can't scan query result: %w",err)
	 }

	 return false,nil
}

func (d MYSQLDB) Register(u entity.User) (entity.User , error){
	res ,err := d.db.Exec(`insert into users(name, phone_number) values(?,?)`, u.Name, u.PhoneNumber)
	if err!=nil{
		return entity.User{},fmt.Errorf("can't execute command: %w",err)
	}

	//error is always nil
	id ,_ := res.LastInsertId()
	u.ID =uint(id)

	return u ,nil
}