package main

import (
	"database/sql"
	"strconv"
	"strings"
)

var defaultConfigs = map[string]string{
	//用于分别使用在Registry和本服务计算Token的证书对
	//该文件内容会自动输出到文件系统中，（必要的话）供Registry调用
	//生成命令参考：
	//  openssl req -new -newkey rsa:4096 -days 365 -subj "/CN=localhost" -nodes -x509 -keyout auth.key -out auth.crt
	"registry_auth_token_key":       "TBD",
	"registry_auth_token_cert":      "TBD",
	"registry_auth_token_key_path":  "/etc/registry_auth_token.key",
	"registry_auth_token_cert_path": "/etc/registry_auth_token.crt",

	//Registry后端服务的地址
	"registry_backend_addr": "http://localhost:5000",

	//docker pull默认采用HTTPS，如果外部没有HTTPS反向代理
	//可以直接把HTTPS证书内容配置在这里，
	//则会按照HTTPS协议监听端口并处理请求
	"registry_https_key":  "",
	"registry_https_cert": "",

	//下面两项和Registry后端保持一致即可
	"registry_token_issuer":       "Issuer",
	"registry_token_service_name": "DockerRegistry",

	//每次请求生成的Token过期时间
	"registry_token_expiration": "1000000000",

	//是否启用Registry仓库的UI API
	"enable_ui": "false",
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

	stmt, err := db.Prepare("INSERT IGNORE INTO `configs` (`name`, `value`) VALUES (?,?)")
	if err != nil {
		return err
	}

	for name, value := range defaultConfigs {
		_, err := stmt.Exec(name, value)
		if err != nil {
			return err
		}
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
