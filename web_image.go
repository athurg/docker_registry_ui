package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"registry"
)

//每一层信息
//镜像commit、命令
func getImageInfo(repo, ref string) (registry.ManifestV2, registry.ImageConfig, error) {
	resourceAction := ResourceActions{Type: "repository", Name: repo, Actions: []string{"pull"}}
	tokenServiceName, _ := GetStringConfig("registry_token_service_name")
	token, err := CreateToken("", tokenServiceName, []ResourceActions{resourceAction})
	if err != nil {
		return registry.ManifestV2{}, registry.ImageConfig{}, fmt.Errorf("创建Token错误: %s", err)
	}

	manifest, err := registryClient.ImageManifestV2(repo, ref, token)
	if err != nil {
		return registry.ManifestV2{}, registry.ImageConfig{}, fmt.Errorf("获取Manifest错误: %s", err)
	}

	config, err := registryClient.ImageConfigByDigest(repo, manifest.Config.Digest, token)
	if err != nil {
		return registry.ManifestV2{}, registry.ImageConfig{}, fmt.Errorf("获取配置错误: %s", err)
	}

	return manifest, config, nil
}

func ApiImageDeleteHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.FormValue("tag")
	repo := r.FormValue("repo")

	tokenServiceName, _ := GetStringConfig("registry_token_service_name")
	resourceAction := ResourceActions{Type: "repository", Name: repo, Actions: []string{"*"}}
	token, err := CreateToken("", tokenServiceName, []ResourceActions{resourceAction})
	if err != nil {
		log.Println(err)
		return
	}

	err = registryClient.ManifestDelete(repo, tag, token)
	if err != nil {
		renderError(w, err)
		return
	}

	renderSuccess(w)
}

func ApiImageShowHandler(w http.ResponseWriter, r *http.Request) {

	image := r.FormValue("name")
	slice := strings.Split(image, ":")
	if len(slice) != 2 {
		renderError(w, fmt.Errorf("镜像名%s不合法", image))
		return
	}

	repo := slice[0]
	ref := slice[1]

	manifest, config, _ := getImageInfo(repo, ref)
	info := map[string]interface{}{
		"Manifest": manifest,
		"Config":   config,
	}

	renderInfo(w, info)
}
