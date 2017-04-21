package main

import (
	"fmt"
	"net/http"
)

func getRepoList() ([]string, error) {
	resourceAction := ResourceActions{Type: "registry", Name: "catalog", Actions: []string{"*"}}
	token, err := CreateToken("", CfgTokenService, []ResourceActions{resourceAction})
	if err != nil {
		return nil, fmt.Errorf("创建Token错误: %s", err)
	}
	err, catalogInfo := registryClient.GetCatalog(token)
	if err != nil {
		return nil, fmt.Errorf("获取仓库列表错误: %s", err)
	}

	return catalogInfo.Repositories, nil
}

func ViewIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := htmlHead

	repos, err := getRepoList()
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
		for _, repo := range repos {
			tags, _ := getRepoTagList(repo)
			html += fmt.Sprintf("<tr><td><a href='/view/repo?name=%s'>%s</a></td><td>%d</td></tr>", repo, repo, len(tags))
		}
		html += `</tbody></table></div>`
	}

	html += htmlFoot
	w.Write([]byte(html))

}
