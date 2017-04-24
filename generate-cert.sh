#!/bin/bash
# 生成Token签名证书
# 生成的.key文件作为KEY_PEM_BLOCK环境变量
# 生成的.cert文件作为CERT_PEM_BLOCK环境变量，并作为Docker Registry的REGISTRY_AUTH_TOKEN_ROOTCERTBUNDLE配置的证书文件
#

if [[ "$1" == ""  ]] ; then
	echo "Usage: $0 basename"
	exit 1
fi

openssl req -new -newkey rsa:4096 -days 365 -subj "/CN=localhost" -nodes -x509 -keyout $1.key -out $1.cert
