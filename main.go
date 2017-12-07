package main

import (
	"log"
	"os"
	"registry"
)

var registryClient *registry.Client

func main() {
	//预解析Registry Token密钥对
	err := preLoadCertAndKey()
	if err != nil {
		log.Fatalf("解析Token签名证书对失败: %s", err)
	}

	//创建Registry后端访问客户端
	registryBackendAddr, _ := GetConfigAsString("registry_backend_addr")
	registryClient = registry.New(registryBackendAddr)

	//创建Registry后端代理及WEB服务
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr != "" {
		listenAddr = ":80"
	}

	err = LoadWebServer(listenAddr, registryBackendAddr)
	if err != nil {
		log.Fatalf("WEB服务启动失败: %s", err)
	}
}
