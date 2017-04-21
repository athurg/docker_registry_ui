package main

import (
	"os"
	"strings"
	"testing"
)

func TestCreateToken(t *testing.T) {
	account := os.Getenv("ACCOUNT")
	service := os.Getenv("SERVICE")
	resourceAction := ResourceActions{
		Type:    os.Getenv("TYPE"),
		Name:    os.Getenv("NAME"),
		Actions: strings.Split(os.Getenv("ACTIONS"), ","),
	}

	token, err := CreateToken(account, service, []ResourceActions{resourceAction})
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("账户:", account)
	t.Log("服务器:", service)
	t.Log("资源类型:", resourceAction.Type)
	t.Log("资源名:", resourceAction.Name)
	t.Log("授权操作:", resourceAction.Actions)
	t.Log("TOKEN:", token)
}
