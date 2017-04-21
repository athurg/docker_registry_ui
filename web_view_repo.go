package main

import (
	"fmt"
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
	resourceAction := ResourceActions{Type: "repository", Name: repo, Actions: []string{"pull"}}
	token, err := CreateToken("", CfgTokenService, []ResourceActions{resourceAction})
	if err != nil {
		return nil, fmt.Errorf("创建Token错误: %s", err)
	}
	err, info := registryClient.GetTags(repo, token)
	if err != nil {
		return nil, fmt.Errorf("获取仓库标签错误: %s", err)
	}

	result := make([]ImageBaseInfo, 0)

	for _, tag := range info.Tags {
		info1, err := registryClient.GetImageManifestV1(repo, tag, token)
		if err != nil {
			return nil, fmt.Errorf("获取仓库标签(%s:%s)详情错误: %s", repo, tag, err)
		}

		info2, err := registryClient.GetImageManifestV2(repo, tag, token)
		if err != nil {
			return nil, fmt.Errorf("获取仓库标签(%s:%s)详情错误: %s", repo, tag, err)
		}

		totalSize := 0
		for _, layer := range info2.Layers {
			totalSize += layer.Size
		}
		topLayer := info1.History[0]
		result = append(result, ImageBaseInfo{
			Id:         info1.History[0].Id,
			Name:       repo + ":" + tag,
			Tag:        tag,
			Repo:       repo,
			Created:    topLayer.Created,
			LayerCount: len(info2.Layers),
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
	<li><a href="/view"><span class="glyphicon glyphicon-calendar"></span>仓库</a></li>
	<li><a href="#">` + repo + `</a></li>
	</ol>
	</div>
	<div class="row">
	<table class="table table-bordered table-hover">
	<thead>
	<tr><th>ID</th><th>标签</th><th>创建时间</th><th>层数</th><th>大小</th></tr>
	</thead>
	<tbody>
	`
	for _, info := range result {
		html += "<tr>"
		html += "<td>" + info.Id[:11] + "</td>"
		html += "<td><a href='/view/image?name=" + info.Name + "'>" + info.Tag + "</a></td>"
		html += "<td>" + info.Created.Format("2006-01-02 15:04:05") + "</td>"
		html += fmt.Sprintf("<td>%d</td>", info.LayerCount)
		html += fmt.Sprintf("<td>%d</td>", info.Size)
		html += "</tr>"
	}
	html += `</tbody></table></div>`

	html += htmlFoot
	w.Write([]byte(html))
}
