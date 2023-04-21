package database

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
)

var Conn *sql.DB

func Open(driverName string, datasourceName string) error {
	var err error
	Conn, err = sql.Open(driverName, datasourceName)
	if err != nil {
		return err
	}
	return nil
}
