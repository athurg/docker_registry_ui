package main

import (
	"fmt"
	"net/http"
	"strings"

	"registry"
)

//每一层信息
//镜像commit、命令
func getImageInfo(repo, ref string) (registry.ManifestV2, registry.ImageConfig, error) {
	resourceAction := ResourceActions{Type: "repository", Name: repo, Actions: []string{"pull"}}
	token, err := CreateToken("", CfgTokenService, []ResourceActions{resourceAction})
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

func DeleteImageHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.FormValue("tag")
	repo := r.FormValue("repo")

	returnUrl := "/view/repo?name=" + repo

	resourceAction := ResourceActions{Type: "repository", Name: repo, Actions: []string{"*"}}
	token, err := CreateToken("", CfgTokenService, []ResourceActions{resourceAction})
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<html><body><h3>请求Token失败%s</h3><div><a href='%s'>返回</a></div></body>", err, returnUrl)
		return
	}

	err = registryClient.ManifestDelete(repo, tag, token)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<html><body><h3>删除镜像失败%s</h3><div><a href='%s'>返回</a></div></body>", err.Error(), returnUrl)
		return
	}

	http.Redirect(w, r, returnUrl, http.StatusMovedPermanently)
}

func ViewImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := htmlHead

	image := r.FormValue("name")
	slice := strings.Split(image, ":")
	if len(slice) != 2 {
		html += "<h3>镜像名" + image + "不合法</h3>" + htmlFoot
		w.Write([]byte(html))
		return
	}

	repo := slice[0]
	ref := slice[1]

	manifest, config, err := getImageInfo(repo, ref)
	if err != nil {
		html += "<h3>" + err.Error() + "</h3>"
		w.Write([]byte(html))
		return
	}

	var tbody string
	var totalSize int

	var layerIdx int

	for _, history := range config.History {
		var size int
		if !history.EmptyLayer {
			size = manifest.Layers[layerIdx].Size
			layerIdx += 1
		}

		totalSize += size

		author := config.Author
		if history.Author != "" {
			author = history.Author
		}

		cmd := history.CreatedBy
		if len(cmd) > 40 {
			cmd = cmd[:40] + "..."
		}

		tbody += "<tr>"
		tbody += fmt.Sprintf("<td>%s</td>", history.Created.Format("2006-01-02 15:04:05"))
		tbody += fmt.Sprintf("<td>%s</td>", cmd)
		tbody += fmt.Sprintf("<td>%s</td>", author)
		tbody += fmt.Sprintf("<td>%s</td>", HumanSize(size))

		tbody += "</tr>"
	}

	html += `
	<div class="row">
	<ol class="breadcrumb">
	<li><a href="/view">仓库</a></li>
	<li><a href="/view/repo?name=` + repo + `">` + repo + `</a></li>
	<li><a href="#">` + ref + `</a></li>
	</ol>
	</div>
	<div class="row">
	<p class="lead">总大小:<mark>` + HumanSize(totalSize) + `</mark></p>
	</div>`

	if len(config.Config.Labels) > 0 {
		html += `
		<div class="row">
			<div class="panel panel-primary">
				<div class="panel-heading">标签</div>
				<div class="panel-body">
					<dl class="dl-horizontal">`
		for label, value := range config.Config.Labels {
			html += `<dt>` + label + `</dt>`
			html += `<dd>` + value + `</dd>`
		}
		html += `
					</dl>
				</div>
			</div><!--panel-->
		</div><!--row-->`
	}

	html += `
	<div class="row">
		<div class="panel panel-primary">
			<div class="panel-heading">历史</div>
			<div class="panel-body">
				<small>
				<table class="table table-bordered table-hover table-condensed">
					<thead>
						<tr>
							<th>创建时间</th>
							<th>命令</th>
							<th>作者</th>
							<th>大小</th>
						</tr>
					</thead>
					<tbody>` + tbody + `</tbody>
				</table>
				</small>
			</div>
		</div><!--panel-->
	</div><!--row-->`

	html += htmlFoot
	w.Write([]byte(html))
}
