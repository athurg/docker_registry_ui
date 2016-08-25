package main

import (
	"os"
	"strconv"
)

var (
	CfgKeyPEMBlock  []byte
	CfgCertPEMBlock []byte

	CfgListenAddr         string = ":8080"
	CfgTokenIssuer        string = "Tap4Fun"
	CfgTokenExpiration    int64  = 10000000
	CfgUserTableName      string = "users"
	CfgDSN                string = "user:password@tcp(host:port)/dbname?charset=utf8&parseTime=True&loc=Local"
	CfgPrivilegeTableName string = "privileges"
)

func init() {
	if v := os.Getenv("KEY_PEM_BLOCK"); v != "" {
		CfgKeyPEMBlock = []byte(v)
	}

	if v := os.Getenv("CERT_PEM_BLOCK"); v != "" {
		CfgCertPEMBlock = []byte(v)
	}

	if v := os.Getenv("DSN"); v != "" {
		CfgDSN = v
	}

	if v := os.Getenv("LISTEN_ADDR"); v != "" {
		CfgListenAddr = v
	}

	if v := os.Getenv("TOKEN_ISSUER"); v != "" {
		CfgTokenIssuer = v
	}

	if v := os.Getenv("USER_TABLE_NAME"); v != "" {
		CfgUserTableName = v
	}

	if v := os.Getenv("PRIVILEGE_TABLE_NAME"); v != "" {
		CfgPrivilegeTableName = v
	}

	if v := os.Getenv("TOKEN_EXPIRATION"); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			CfgTokenExpiration = int64(n)
		}
	}
}
