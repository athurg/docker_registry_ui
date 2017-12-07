package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	KiB = 1024
	MiB = 1024 * KiB
	GiB = 1024 * MiB
)

func HumanSize(size int) string {
	switch {
	case size > GiB:
		return fmt.Sprintf("%.2f GiB", float32(size)/float32(GiB))
	case size > MiB:
		return fmt.Sprintf("%.2f MiB", float32(size)/float32(MiB))
	case size > KiB:
		return fmt.Sprintf("%.2f KiB", float32(size)/float32(KiB))
	default:
		return fmt.Sprintf("%.0f", float32(size))
	}
}

const (
	registryHttpsKeyFile  = "/etc/registry_https.key"
	registryHttpsCertFile = "/etc/registry_https.crt"
)

func LoadWebServer(addr, registryBackendAddr string) error {
	registryBackendURL, err := url.Parse(registryBackendAddr)
	if err != nil {
		return fmt.Errorf("Docker Registry 地址解析失败: %s", err)
	}

	//对Registry对请求作代理
	registryProxy := httputil.NewSingleHostReverseProxy(registryBackendURL)
	http.HandleFunc("/v2/", registryProxy.ServeHTTP)
	http.HandleFunc("/v1/", registryProxy.ServeHTTP)

	//其他请求自行处理
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/auth", AuthHandler)
	http.HandleFunc("/view", ViewIndexHandler)
	http.HandleFunc("/view.json", ViewIndexJSONHandler)
	http.HandleFunc("/view/repo", ViewRepoHandler)
	http.HandleFunc("/view/image", ViewImageHandler)
	http.HandleFunc("/view/image/delete", DeleteImageHandler)

	//如果提供了HTTPS的密钥对，则监听为HTTPS，否则监听为HTTP
	registryHttpsKeyBlock, _ := GetConfigAsString("registry_https_key")
	registryHttpsCertBlock, _ := GetConfigAsString("registry_https_cert")
	if registryHttpsKeyBlock == "" || registryHttpsCertBlock == "" {
		return http.ListenAndServe(addr, nil)
	}

	ioutil.WriteFile(registryHttpsKeyFile, []byte(registryHttpsKeyBlock), 0755)
	ioutil.WriteFile(registryHttpsCertFile, []byte(registryHttpsCertBlock), 0755)

	return http.ListenAndServeTLS(addr, registryHttpsCertFile, registryHttpsKeyFile, nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("请求来了", r.URL)
	http.Redirect(w, r, "/view", http.StatusMovedPermanently)
}
