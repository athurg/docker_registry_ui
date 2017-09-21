package main

import (
	"log"
	"registry"
)

var registryClient *registry.Client

func main() {
	err := ParseConfig()
	if err != nil {
		log.Fatalf("解析配置失败: %s", err)
	}

	err = LoadCertAndKey()
	if err != nil {
		log.Fatalf("解析Token签名证书对失败: %s", err)
	}

	registryClient = registry.New(CfgRegistryBackendAddr)
	err = LoadWebServer(CfgListenAddr, CfgRegistryBackendAddr)
	if err != nil {
		log.Fatalf("WEB服务启动失败: %s", err)
	}
}
