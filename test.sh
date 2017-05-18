#!/bin/sh

#线上模式
export ACCOUNT="sparta"
export TOKEN_ISSUER="SpartaIdp"
export REGISTRY_ADDR="http://172.20.100.189:5000"
export TOKEN_SERVICE_NAME="DockerRegistry"
export KEY_PEM_BLOCK=$(cat key.pem)
export CERT_PEM_BLOCK=$(cat crt.pem)

GOPATH=`pwd` go install app && ./bin/app
