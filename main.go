package main

import (
	"log"
)

func init() {
	err := LoadCertAndKey()
	if err != nil {
		log.Fatalf("Failed to parse cert and key: %s", err)
	}
}

func main() {
	LoadWebServer(CfgListenAddr)
}
