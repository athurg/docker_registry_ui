package main

import (
	"log"
	"registry"
)

var registryClient *registry.Client

func init() {
	err := LoadCertAndKey()
	if err != nil {
		log.Fatalf("Failed to parse cert and key: %s", err)
	}

	registryClient = registry.New(CfgRegistryAddr)
}

func main() {
	LoadWebServer(CfgListenAddr)
}
