package main

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Config struct {
	Key   string
	Value string
}

func GetConfigAsString(key string) (string, error) {
	var value string
	row := dbConn.QueryRow("SELECT `value` FROM `configs` WHERE `key`=?", key)
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("User not exists")
	}

	return value, nil
}

func GetConfigAsInt64(key string) (int64, error) {
	var value string
	row := dbConn.QueryRow("SELECT `value` FROM `configs` WHERE `key`=?", key)
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("User not exists")
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return int64(v), nil
}
