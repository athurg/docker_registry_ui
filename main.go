package main

import (
	"log"
	"os"
	"registry"
)

var registryClient *registry.Client

func main() {
	//创建Registry后端访问客户端
	registryBackendAddr, _ := GetStringConfig("registry_backend_addr")
	registryClient = registry.New(registryBackendAddr)

	//创建Registry后端代理及WEB服务
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":80"
	}

	err := LoadWebServer(listenAddr, registryBackendAddr)
	if err != nil {
		log.Fatalf("WEB服务启动失败: %s", err)
	}
}
