package main

import (
	"fmt"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB

func connectDb() error {
	if dbConn == nil || dbConn.Ping() != nil {
		db, err := sql.Open("mysql", CfgDSN)
		if err != nil {
			return fmt.Errorf("DB connect failed: %s", err)
		}
		dbConn = db
	}
	return nil
}
