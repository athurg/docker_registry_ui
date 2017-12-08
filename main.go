package main

import (
	"log"
	"registry"
)

var registryClient *registry.Client

func main() {
	//创建Registry后端访问客户端
	registryBackendAddr, _ := GetStringConfig("registry_backend_addr")
	registryClient = registry.New(registryBackendAddr)

	err := LoadWebServer(registryBackendAddr)
	if err != nil {
		log.Fatalf("WEB服务启动失败: %s", err)
	}
}
