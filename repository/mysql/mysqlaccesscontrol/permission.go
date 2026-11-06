package mysqlaccesscontrol

import (
	"gameApp/entity"
	"gameApp/repository/mysql"
	"time"
)

func scanPermission(scanner mysql.Scanner) (entity.Permission, error) {
	var createdAt time.Time
	var p entity.Permission
	//It must be in the exact same order as the columns defined in the database.
	//We also need to pass its pointer so that the scan function can set the values.
	err := scanner.Scan(&p.ID, &p.Title, &createdAt)

	return p, err
}