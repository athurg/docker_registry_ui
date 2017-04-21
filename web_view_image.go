package main

import (
	"fmt"
	"net/http"
	"strings"

	"./registry"
)

//每一层信息
//镜像commit、命令

func getImageManifest(repo, ref string) (registry.ImageManifest, registry.ImageManifest, error) {
	resourceAction := ResourceActions{Type: "repository", Name: repo, Actions: []string{"pull"}}
	token, err := CreateToken("", CfgTokenService, []ResourceActions{resourceAction})
	if err != nil {
		return registry.ImageManifest{}, registry.ImageManifest{}, fmt.Errorf("创建Token错误: %s", err)
	}

	info2, err := registryClient.GetImageManifestV2(repo, ref, token)
	if err != nil {
		return registry.ImageManifest{}, registry.ImageManifest{}, fmt.Errorf("获取仓库标签错误: %s", err)
	}

	info, err := registryClient.GetImageManifestV1(repo, ref, token)
	if err != nil {
		return registry.ImageManifest{}, registry.ImageManifest{}, fmt.Errorf("获取仓库容器配置标签错误: %s", err)
	}

	return info, info2, nil
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

	info, info2, err := getImageManifest(repo, ref)
	if err != nil {
		html += "<h3>" + err.Error() + "</h3>"
		w.Write([]byte(html))
		return
	}

	//
	sizeByDigest := make(map[string]int)
	for _, layer := range info2.Layers {
		sizeByDigest[layer.Digest] = layer.Size
	}

	var tbody string
	var totalSize int
	for i, history := range info.History {
		tbody += "<tr>"
		tbody += fmt.Sprintf("<td>%s</td>", history.Id[:11])
		size := sizeByDigest[info.FsLayers[i].BlobSum]
		tbody += fmt.Sprintf("<td>%d</td>", size)
		tbody += fmt.Sprintf("<td>%s</td>", history.Author)
		tbody += fmt.Sprintf("<td>%s</td>", history.Created.Format("2006-01-02 15:04:05"))

		cmd := strings.Join(history.ContainerConfig.Cmd, "")
		if len(cmd) > 40 {
			cmd = cmd[:40] + "..."
		}
		tbody += fmt.Sprintf("<td>%s</td>", cmd)

		totalSize += sizeByDigest[info.FsLayers[i].BlobSum]
		tbody += "</tr>"
	}

	html += `
	<div class="row">
	<ol class="breadcrumb">
	<li><a href="/view"><span><span class="glyphicon glyphicon-calendar">仓库</a></li>
	<li><a href="/view/repo?name=` + repo + `">` + repo + `</a></li>
	<li><a href="#">` + ref + `</a></li>
	</ol>
	</div>
	<div class="row">
	<p>总大小:` + fmt.Sprintf("%d", totalSize) + `</p>
	<table class="table table-bordered table-hover">
	<thead>
	<tr><th>ID</th><th>大小</th><th>作者</th><th>创建时间</th><th>命令</th></tr>
	</thead>
	<tbody>` + tbody + `</tbody></table></div>`

	html += htmlFoot
	w.Write([]byte(html))
}
