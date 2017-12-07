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

func ViewRepoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := htmlHead

	repo := r.FormValue("name")
	if repo == "" {
		html += "<h3>未指定仓库名</h3>" + htmlFoot
		w.Write([]byte(html))
		return
	}

	result, err := getRepoTagList(repo)
	if err != nil {
		html += "<h3>" + err.Error() + "</h3>"
		w.Write([]byte(html))
		return
	}

	html += `
	<div class="row">
	<ol class="breadcrumb">
	<li><a href="/view">仓库</a></li>
	<li><a href="#">` + repo + `</a></li>
	</ol>
	</div>
	<div class="row">
	<table class="table table-bordered table-hover">
	<thead>
	<tr><th>ID</th><th>标签</th><th>创建时间</th><th>层数</th><th>大小</th><th>操作</th></tr>
	</thead>
	<tbody>
	`
	for _, info := range result {
		html += "<tr>"
		html += "<td>" + info.Id[:12] + "</td>"
		html += "<td><a href='/view/image?name=" + info.Name + "'>" + info.Tag + "</a></td>"
		html += "<td>" + info.Created.Format("2006-01-02 15:04:05") + "</td>"
		html += fmt.Sprintf("<td>%d</td>", info.LayerCount)
		html += fmt.Sprintf("<td>%s</td>", HumanSize(info.Size))
		html += fmt.Sprintf("<td><a href='/view/image/delete?repo=%s&tag=%s'>删除</td>", info.Repo, info.Tag)
		html += "</tr>"
	}
	html += `</tbody></table></div>`

	html += htmlFoot
	w.Write([]byte(html))
}
