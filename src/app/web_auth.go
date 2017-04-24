package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sort"
	"strings"
)

type AuthRequest struct {
	RemoteIP net.IP
	Password string
	Account  string
	Service  string
	Scopes   []AuthScope
}

func parseRemoteAddr(ra string) net.IP {
	colonIndex := strings.LastIndex(ra, ":")
	if colonIndex > 0 && ra[colonIndex-1] >= 0x30 && ra[colonIndex-1] <= 0x39 {
		ra = ra[:colonIndex]
	}
	if ra[0] == '[' && ra[len(ra)-1] == ']' { // IPv6
		ra = ra[1 : len(ra)-1]
	}
	res := net.ParseIP(ra)
	return res
}

func parseRequest(req *http.Request) (*AuthRequest, error) {
	if err := req.ParseForm(); err != nil {
		return nil, fmt.Errorf("invalid form value")
	}

	if realAddr := req.Header.Get("X-Forwarded-For"); realAddr != "" {
		req.RemoteAddr = realAddr
	}

	ip := parseRemoteAddr(req.RemoteAddr)
	if ip == nil {
		return nil, fmt.Errorf("unable to parse remote addr %s", req.RemoteAddr)
	}

	authReq := &AuthRequest{RemoteIP: ip}

	user, pass, basicAuth := req.BasicAuth()
	if basicAuth {
		authReq.Account = user
		authReq.Password = pass
	}

	account := req.FormValue("account")
	if basicAuth && account != authReq.Account {
		return nil, fmt.Errorf("user and account are not the same (%q vs %q)", authReq.Account, account)
	}

	//TODO: 也许可以尝试校验Service
	authReq.Service = req.FormValue("service")

	for _, scopeStr := range req.Form["scope"] {
		parts := strings.Split(scopeStr, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid scope: %q", scopeStr)
		}

		scope := AuthScope{
			RepoName: parts[1],
			Category: parts[0],
			Actions:  strings.Split(parts[2], ","),
		}
		sort.Strings(scope.Actions)
		authReq.Scopes = append(authReq.Scopes, scope)
	}
	return authReq, nil
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	ar, err := parseRequest(r)
	if err != nil {
		log.Printf("Bad request: %s", err)
		http.Error(w, fmt.Sprintf("Bad request: %s", err), http.StatusBadRequest)
		return
	}

	log.Printf("[INFO]%s请求授权服务%s, %s", ar.Account, ar.Service, ar.Scopes)

	if len(ar.Account) == 1 {
		log.Printf("用户名长度必须大于1(非匿名用户)或等于0(匿名用户): %s", ar.Account)
		http.Error(w, fmt.Sprintf("Invalid username %s, length must large than 1", ar.Account), http.StatusUnauthorized)
		return
	}

	//用户鉴定
	if ar.Account == "" {
		ar.Account = "*"
	}
	u, err := GetUser(ar.Account, ar.Password)
	if err != nil {
		log.Printf("用户查找失败: %s", err)
		http.Error(w, fmt.Sprintf("Invalid user: %s", err), http.StatusUnauthorized)
		return
	}

	//用户鉴权
	authzResults := []ResourceActions{}
	if len(ar.Scopes) > 0 {
		authzResults, err = u.Authorize(ar.RemoteIP, ar.Scopes)
		if err != nil {
			http.Error(w, fmt.Sprintf("Authorization failed (%s)", err), http.StatusInternalServerError)
			return
		}
	} else {
		// Authentication-only request ("docker login"), pass through.
	}

	token, err := CreateToken(ar.Account, ar.Service, authzResults)
	if err != nil {
		msg := fmt.Sprintf("Failed to generate token %s", err)
		http.Error(w, msg, http.StatusInternalServerError)
		log.Printf("%s: %s", ar, msg)
		return
	}
	result, _ := json.Marshal(&map[string]string{"token": token})

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
