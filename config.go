package main

import (
	"fmt"
	"os"
	"strconv"
)

var (
	//WEB服务监听地址（视HTTPS证书提供与否，有可能是HTTP或HTTPS服务）
	CfgListenAddr string = ":80"

	CfgRegistryBackendAddr string = "http://localhost:5000"

	CfgDSN string = "user:password@tcp(host:port)/dbname?charset=utf8&parseTime=True&loc=Local"

	//签名的相关信息，应该和Registry后端中保持一致
	CfgTokenIssuer     string = "Issurer"
	CfgTokenService    string = "Registry Service"
	CfgTokenExpiration int64  = 10000000
)

func ParseConfig() error {
	if v := os.Getenv("DSN"); v != "" {
		CfgDSN = v
	}

	if v := os.Getenv("REGISTRY_BACKEND_ADDR"); v != "" {
		CfgRegistryBackendAddr = v
	}

	if v := os.Getenv("REGISTRY_AUTH_TOKEN_SERVICE"); v != "" {
		CfgTokenService = v
	}

	if v := os.Getenv("LISTEN_ADDR"); v != "" {
		CfgListenAddr = v
	}

	if v := os.Getenv("REGISTRY_AUTH_TOKEN_ISSUER"); v != "" {
		CfgTokenIssuer = v
	}

	if v := os.Getenv("TOKEN_EXPIRATION"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("无效的TOKEN有效期: %s", v)
		}
		CfgTokenExpiration = int64(n)
	}

	return nil
}
