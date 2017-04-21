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

func getRepoTagList(repo string) ([]string, error) {
	resourceAction := ResourceActions{Type: "repository", Name: repo, Actions: []string{"pull"}}
	token, err := CreateToken("", CfgTokenService, []ResourceActions{resourceAction})
	if err != nil {
		return nil, fmt.Errorf("创建Token错误: %s", err)
	}
	err, info := registryClient.GetTags(repo, token)
	if err != nil {
		return nil, fmt.Errorf("获取仓库标签错误: %s", err)
	}

	return info.Tags, nil
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<html lang="zh-CN">
	<head>
		<title>Docker镜像浏览器</title>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
		<script src="https://cdn.bootcss.com/jquery/1.12.4/jquery.min.js"></script>
		<script src="https://cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>
	</head>
	<body>
		<div class="container">
	`

	repo := r.FormValue("repo")
	if repo != "" {
		tags, err := getRepoTagList(repo)
		if err != nil {
			html += "<p>" + err.Error() + "</p>"
		} else {
			html += `
			<div class="row">
				<ol class="breadcrumb">
					<li><a href="view"><span><span class="glyphicon glyphicon-calendar"></span>仓库</span></a></li>
					<li><a href="#"><span><span class="glyphicon glyphicon-calendar"></span>` + repo + `</span></a></li>
				</ol>
			</div>
			<div class="row">
				<table class="table table-bordered table-hover"><thead><tr><th>标签</th></tr></thead><tbody>
			`
			for _, tag := range tags {
				imageName := repo + ":" + tag
				html += "<tr><td><a href='?image=" + imageName + "'>" + imageName + "</a></td></tr>"
			}
			html += `</tbody></table></div>`
		}
	} else {
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
				html += fmt.Sprintf("<tr><td><a href='?repo=%s'>%s</a></td><td>%d</td></tr>", repo, repo, len(tags))
			}
			html += `</tbody></table></div>`
		}
	}

	html += `
	</div>
	</body>
	</html>`
	w.Write([]byte(html))

}
