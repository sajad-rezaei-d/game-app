package mysql

import (
	"database/sql"
	"fmt"
	"gameApp/entity"
	"gameApp/pkg/errmsg"
	"gameApp/pkg/richerror"
	"time"
)

func (d *MYSQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	row := d.db.QueryRow(`select * from users where phone_number = ?`, phoneNumber)

	_, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, fmt.Errorf("can't scan query result: %w", err)
	}

	return false, nil
}

func (d *MYSQLDB) Register(u entity.User) (entity.User, error) {
	res, err := d.db.Exec(`insert into users(name, phone_number, password) values(?,?,?)`, u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("can't execute command: %w", err)
	}

	//error is always nil
	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d *MYSQLDB) GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error) {
	const op = "mysql.GetUserByPhoneNumber"

	row := d.db.QueryRow(`select * from users where phone_number = ?`, phoneNumber)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, false, nil
		}
		// fmt.Errorf("can't scan query result: %w",err)
		return entity.User{}, false,richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryREsult).WithKind(richerror.KindUnexpected)
	}

	return user, true, nil
}

func (d *MYSQLDB) GetUserByID(userID uint) (entity.User, error) {
	const op = "mysql.GetUserByID"
	row := d.db.QueryRow(`select * from users where id = ?`, userID)
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

func scanUser(row *sql.Row) (entity.User, error) {
	var createdAt time.Time
	var user entity.User
	//It must be in the exact same order as the columns defined in the database.
	//We also need to pass its pointer so that the scan function can set the values.
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt, &user.Password)

	return user, err
}
