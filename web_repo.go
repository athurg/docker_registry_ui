package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type ImageBaseInfo struct {
	Id         string
	Name       string
	Repo       string
	Tag        string
	Created    time.Time
	LayerCount int
	Size       int
}

func getRepoTagList(repo string) ([]ImageBaseInfo, error) {
	tokenServiceName, _ := GetConfigAsString("registry_token_service_name")
	resourceAction := ResourceActions{Type: "repository", Name: repo, Actions: []string{"pull"}}
	token, err := CreateToken("", tokenServiceName, []ResourceActions{resourceAction})
	if err != nil {
		return nil, fmt.Errorf("创建Token错误: %s", err)
	}
	err, info := registryClient.GetTags(repo, token)
	if err != nil {
		return nil, fmt.Errorf("获取仓库标签错误: %s", err)
	}

	result := make([]ImageBaseInfo, 0)

	for _, tag := range info.Tags {
		manifest, err := registryClient.ImageManifestV2(repo, tag, token)
		if err != nil {
			log.Printf("获取仓库标签(%s:%s)详情错误: %s", repo, tag, err)
			continue
		}

		config, err := registryClient.ImageConfigByDigest(repo, manifest.Config.Digest, token)
		if err != nil {
			log.Printf("获取仓库标签(%s:%s)详情错误: %s", repo, tag, err)
			continue
		}

		totalSize := 0
		for _, layer := range manifest.Layers {
			totalSize += layer.Size
		}

		result = append(result, ImageBaseInfo{
			Id:         manifest.Config.Digest[7:],
			Name:       repo + ":" + tag,
			Tag:        tag,
			Repo:       repo,
			Created:    config.Created,
			LayerCount: len(manifest.Layers),
			Size:       totalSize,
		})
	}

	return result, nil
}

func ApiRepoShowHandler(w http.ResponseWriter, r *http.Request) {
	repo := r.FormValue("name")
	if repo == "" {
		renderError(w, fmt.Errorf("未指定仓库名"))
		return
	}

	result, err := getRepoTagList(repo)
	if err != nil {
		renderError(w, err)
		return
	}

	renderInfo(w, result)
}

func getRepoList() ([]string, []int, error) {
	tokenServiceName, _ := GetConfigAsString("registry_token_service_name")
	resourceAction := ResourceActions{Type: "registry", Name: "catalog", Actions: []string{"*"}}
	token, err := CreateToken("", tokenServiceName, []ResourceActions{resourceAction})
	if err != nil {
		return nil, nil, fmt.Errorf("创建Token错误: %s", err)
	}

	err, catalogInfo := registryClient.GetCatalog(token)
	if err != nil {
		return nil, nil, fmt.Errorf("获取仓库列表错误: %s", err)
	}

	tagCounts := make([]int, len(catalogInfo.Repositories))
	for i, repo := range catalogInfo.Repositories {
		tokenServiceName, _ := GetConfigAsString("registry_token_service_name")
		resourceAction := ResourceActions{Type: "repository", Name: repo, Actions: []string{"pull"}}
		token, err := CreateToken("", tokenServiceName, []ResourceActions{resourceAction})
		if err != nil {
			log.Printf("创建Token错误: %s", err)
			continue
		}
		err, info := registryClient.GetTags(repo, token)
		if err != nil {
			log.Printf("获取仓库标签错误: %s", err)
			continue
		}

		tagCounts[i] = len(info.Tags)
	}

	return catalogInfo.Repositories, tagCounts, nil
}

func ApiRepoIndexHandler(w http.ResponseWriter, r *http.Request) {
	repos, tagCount, _ := getRepoList()

	info := map[string]interface{}{
		"TagCount": tagCount,
		"Repos":    repos,
	}

	renderInfo(w, info)
}
