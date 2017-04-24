package main

import (
	"fmt"
	"net/http"
)

func getRepoList() ([]string, []int, error) {
	resourceAction := ResourceActions{Type: "registry", Name: "catalog", Actions: []string{"*"}}
	token, err := CreateToken("", CfgTokenService, []ResourceActions{resourceAction})
	if err != nil {
		return nil, nil, fmt.Errorf("创建Token错误: %s", err)
	}
	err, catalogInfo := registryClient.GetCatalog(token)
	if err != nil {
		return nil, nil, fmt.Errorf("获取仓库列表错误: %s", err)
	}

	tagCounts := make([]int, len(catalogInfo.Repositories))
	for i, repo := range catalogInfo.Repositories {
		resourceAction := ResourceActions{Type: "repository", Name: repo, Actions: []string{"pull"}}
		token, err := CreateToken("", CfgTokenService, []ResourceActions{resourceAction})
		if err != nil {
			return nil, nil, fmt.Errorf("创建Token错误: %s", err)
		}
		err, info := registryClient.GetTags(repo, token)
		if err != nil {
			return nil, nil, fmt.Errorf("获取仓库标签错误: %s", err)
		}

		tagCounts[i] = len(info.Tags)
	}

	return catalogInfo.Repositories, tagCounts, nil
}

func ViewIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := htmlHead

	repos, tagCounts, err := getRepoList()
	if err != nil {
		html += "<p>" + err.Error() + "</p>"
	} else {
		html += `
		<div class="row">
		<ol class="breadcrumb">
		<li><a href="view"><span><span class="glyphicon glyphicon-calendar"></span>仓库</span></a></li>
		</ol>
		</div>
		<div class="row">
		<table class="table table-bordered table-hover"><thead><tr><th>仓库</th><th>镜像数量</th></tr></thead><tbody>
		`
		for i, repo := range repos {
			html += fmt.Sprintf("<tr><td><a href='/view/repo?name=%s'>%s</a></td><td>%d</td></tr>", repo, repo, tagCounts[i])
		}
		html += `</tbody></table></div>`
	}

	html += htmlFoot
	w.Write([]byte(html))

}
