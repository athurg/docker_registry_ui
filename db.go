package main

import (
	"log"
	"os"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB

func init() {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "root:root@tcp(localhost:3306)/dohub?charset=utf8&parseTime=True&loc=Local"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Invalid DSN %s: %s", dsn, err)
		return
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Fail to connect database: %s", err)
		return
	}

	//创建三个表
	if _, err = db.Exec(createUserSql); err != nil {
		log.Fatalf("Fail to create users table: %s", err)
	}
	if _, err = db.Exec(createConfigSql); err != nil {
		log.Fatalf("Fail to create configs table: %s", err)
	}
	if _, err = db.Exec(createPrivilegeSql); err != nil {
		log.Fatalf("Fail to create privileges table: %s", err)
	}

	dbConn = db
}
