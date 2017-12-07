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
		log.Fatal("DB connect failed: %s", err)
		return
	}

	dbConn = db
}
