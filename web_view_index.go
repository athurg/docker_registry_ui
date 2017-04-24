package main

import (
	"fmt"
	"net/http"
)

func getRepoList() (map[string]int, error) {
	resourceAction := ResourceActions{Type: "registry", Name: "catalog", Actions: []string{"*"}}
	token, err := CreateToken("", CfgTokenService, []ResourceActions{resourceAction})
	if err != nil {
		return nil, fmt.Errorf("创建Token错误: %s", err)
	}
	err, catalogInfo := registryClient.GetCatalog(token)
	if err != nil {
		return nil, fmt.Errorf("获取仓库列表错误: %s", err)
	}

	result := make(map[string]int)
	for _, repo := range catalogInfo.Repositories {
		resourceAction := ResourceActions{Type: "repository", Name: repo, Actions: []string{"pull"}}
		token, err := CreateToken("", CfgTokenService, []ResourceActions{resourceAction})
		if err != nil {
			return nil, fmt.Errorf("创建Token错误: %s", err)
		}
		err, info := registryClient.GetTags(repo, token)
		if err != nil {
			return nil, fmt.Errorf("获取仓库标签错误: %s", err)
		}

		result[repo] = len(info.Tags)
	}

	return result, nil
}

func ViewIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := htmlHead

	repoTagCount, err := getRepoList()
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
		for repo, tagCount := range repoTagCount {
			html += fmt.Sprintf("<tr><td><a href='/view/repo?name=%s'>%s</a></td><td>%d</td></tr>", repo, repo, tagCount)
		}
		html += `</tbody></table></div>`
	}

	html += htmlFoot
	w.Write([]byte(html))

}
