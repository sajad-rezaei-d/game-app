package main

import (
	"fmt"
	"gameApp/entity"
	"gameApp/repository/mysql"
)

func main(){
	

}

func TestMysqlRepo(){
	mysqlRepo:=mysql.New()

	createdUser, err:= mysqlRepo.Register(entity.User{
		ID: 0,
		PhoneNumber: "0999",
		Name: "amir",
	})

	if err!=nil {
		fmt.Println("register user",err)
	}else{
		fmt.Println("created user", createdUser)
	}
	isUnique ,err := mysqlRepo.IsPhoneNumberUnique(createdUser.PhoneNumber+"99")
	if err!=nil {
			fmt.Println("unique err",err)
	}

	fmt.Println("isUnique", isUnique)
}