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

	//初始化数据
	if err := InitUserTable(db); err != nil {
		log.Fatalf("Fail to init users table", err)
	}

	if err := InitPrivilegeTable(db); err != nil {
		log.Fatalf("Fail to init users table", err)
	}

	if err := InitConfigTable(db); err != nil {
		log.Fatalf("Fail to init users table", err)
	}

	dbConn = db
}
