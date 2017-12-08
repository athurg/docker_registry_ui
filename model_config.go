package main

import (
	"database/sql"
	"strconv"
	"strings"
)

var configNames = []string{
	"registry_auth_token_key",
	"registry_auth_token_cert",
	"registry_backend_addr",
	"registry_https_key",
	"registry_https_cert",
	"registry_token_issuer",
	"registry_token_service_name",
	"registry_token_expiration",
	"enable_ui",
}

func InitConfigTable(db *sql.DB) error {
	createSql := "CREATE TABLE IF NOT EXISTS `configs` ("
	createSql += "  `name` VARCHAR(255) NOT NULL,"
	createSql += "  `value` text NOT NULL,"
	createSql += "  PRIMARY KEY (`name`)"
	createSql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8"
	if _, err := db.Exec(createSql); err != nil {
		return err
	}

	values := make([]string, 0, len(configNames))
	for _, name := range configNames {
		values = append(values, "('"+name+"', 'TBD')")
	}

	initSql := "INSERT IGNORE INTO `configs` (`name`, `value`) VALUES " + strings.Join(values, ",")
	if _, err := db.Exec(initSql); err != nil {
		return err
	}

	return nil
}

func GetStringConfig(name string) (string, error) {
	var value string
	row := dbConn.QueryRow("SELECT `value` FROM `configs` WHERE `name`=?", name)
	err := row.Scan(&value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func GetInt64Config(name string) (int64, error) {
	var value string
	row := dbConn.QueryRow("SELECT `value` FROM `configs` WHERE `name`=?", name)
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		return 0, err
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return int64(v), nil
}

func GetBoolConfig(name string) (bool, error) {
	var value string
	row := dbConn.QueryRow("SELECT `value` FROM `configs` WHERE `name`=?", name)
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		return false, err
	}

	return strings.ToLower(value) == "true", nil
}
