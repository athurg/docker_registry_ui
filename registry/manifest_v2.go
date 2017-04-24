package registry

import (
	"net/http"
	"time"
)

//镜像Manifest(V 2, Schema 2)
//定义参考: https://docs.docker.com/registry/spec/manifest-v2-2/
type ManifestV2 struct {
	BaseResponse
	SchemaVersion int //始终是2

	//以下为SchemaVersion=2时才有的字段
	MediaType string
	//包含一个用于初始化容器的、JSON格式的、Blob对象
	//可以通过Digest字段获取对应的Blob数据，其内容为JSON格式的配置（猜想）
	Config struct {
		MediaType string //始终是application/vnd.docker.container.image.v1+json
		Size      int    //容器配置对象的大小
		Digest    string //容器配置对象的Digest，可以用这个Digest去pull对应的Blob获取配置的内容
	}
	//从基础镜像开始的层
	Layers []struct {
		MediaType string
		Size      int
		Digest    string
		Urls      interface{}
	}
}

//获取Schema V2格式的Manifest
func (cli *Client) ImageManifestV2(repo, ref, token string) (ManifestV2, error) {
	var info ManifestV2

	header := http.Header{}
	header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	err := cli.RequestWithHeader("/v2/"+repo+"/manifests/"+ref, token, header, &info)

	return info, err
}
