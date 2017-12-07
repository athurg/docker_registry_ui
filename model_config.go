package main

import (
	"database/sql"
	"fmt"
)

type Config struct {
	Key   string
	Value string
}

func GetConfigAsString(key string) (string, error) {
	if err := connectDb(); err != nil {
		return "", fmt.Errorf("无法链接数据库: %s", err)
	}

	var value string
	row := dbConn.QueryRow("SELECT `value` FROM `configs` WHERE `key`=?", key)
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("User not exists")
	}

	return value, nil
}
