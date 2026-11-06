package mysqluser

import (
	"database/sql"
	"fmt"
	"gameApp/entity"
	"gameApp/pkg/errmsg"
	"gameApp/pkg/richerror"
	"gameApp/repository/mysql"
	"time"
)

func (d *DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"
	row := d.conn.Conn().QueryRow(`select * from users where phone_number = ?`, phoneNumber)

	_, err := scanUser(row)
	if err != nil {
		fmt.Println("mysql error",err)
		if err == sql.ErrNoRows {
			fmt.Println("here")
			return true, nil
		}
		return false, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryREsult).WithKind(richerror.KindUnexpected)
	}

	return false, nil
}

func (d *DB) Register(u entity.User) (entity.User, error) {
	res, err := d.conn.Conn().Exec(`insert into users(name, phone_number, password,role) values(?,?,?,?)`, u.Name, u.PhoneNumber, u.Password, u.Role.String())
	if err != nil {
		return entity.User{}, fmt.Errorf("can't execute command: %w", err)
	}

	//error is always nil
	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d *DB) GetUserByPhoneNumber(phoneNumber string) (entity.User,error) {
	const op = "mysql.GetUserByPhoneNumber"

	row := d.conn.Conn().QueryRow(`select * from users where phone_number = ?`, phoneNumber)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
		}
		// fmt.Errorf("can't scan query result: %w",err)
		// TODO: log unexpected error for better observability
		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryREsult).WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d *DB) GetUserByID(userID uint) (entity.User, error) {
	const op = "mysql.GetUserByID"
	row := d.conn.Conn().QueryRow(`select * from users where id = ?`, userID)
	//It must be in the exact same order as the columns defined in the database.
	//We also need to pass its pointer so that the scan function can set the values.
	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
			
		}
		return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgCantScanQueryREsult).WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func scanUser(scanner mysql.Scanner) (entity.User, error) {
	var createdAt time.Time
	var user entity.User

	var roleStr string
	//It must be in the exact same order as the columns defined in the database.
	//We also need to pass its pointer so that the scan function can set the values.
	err := scanner.Scan(&user.ID, &user.Name,&createdAt,&user.PhoneNumber, &user.Password, &roleStr)

	user.Role = entity.MapToRoleEntity(roleStr)

	return user, err
}

