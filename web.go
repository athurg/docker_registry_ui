package main

import (
	"net/http"
)

func LoadWebServer(addr string) {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/auth", AuthHandler)
	http.HandleFunc("/view", ViewIndexHandler)
	http.HandleFunc("/view/repo", ViewRepoHandler)
	http.HandleFunc("/view/image", ViewImageHandler)

	http.ListenAndServe(addr, nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("Hallo"))
}
