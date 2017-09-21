package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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
	http.HandleFunc("/view/repo", ViewRepoHandler)
	http.HandleFunc("/view/image", ViewImageHandler)
	http.HandleFunc("/view/image/delete", DeleteImageHandler)

	httpsCertFile := os.Getenv("REGISTRY_UI_HTTPS_CERT")
	httpsKeyFile := os.Getenv("REGISTRY_UI_HTTPS_KEY")

	if httpsKeyFile == "" || httpsCertFile == "" {
		return http.ListenAndServe(addr, nil)
	} else {
		return http.ListenAndServeTLS(addr, httpsCertFile, httpsKeyFile, nil)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("请求来了", r.URL)
	http.Redirect(w, r, "/view", http.StatusMovedPermanently)
}
