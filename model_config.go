package main

import (
	"database/sql"
	"fmt"
	"strconv"
)

const createConfigSql = "" +
	"CREATE TABLE IF NOT EXISTS `configs` (" +
	"  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT," +
	"  `key` VARCHAR(255) NOT NULL," +
	"  `value` text NOT NULL," +
	"  PRIMARY KEY (`id`)" +
	") ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8"

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
