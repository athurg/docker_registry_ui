package registry

import (
	"encoding/json"
	"net/http"
	"time"
)

type ContainerConfig struct {
	Hostname   string
	Domainname string
	User       string
	Cmd        []string
	Env        []string
	WorkingDir string
	Labels     map[string]string
	Entrypoint []string
	Image      string

	AttachStdin  bool
	OpenStdin    bool
	ArgsEscaped  bool
	AttachStdout bool
	AttachStderr bool
	Tty          bool
	StdinOnce    bool

	Volumes interface{}
	OnBuild []interface{}
}

type History struct {
	Architecture    string
	Container       string
	Os              string
	Parent          string
	Author          string
	Docker_version  string
	Id              string
	Throwaway       bool
	Created         time.Time
	Config          *ContainerConfig
	ContainerConfig *ContainerConfig `json:"container_config"`
}

type ImageManifest struct {
	BaseResponse

	SchemaVersion int

	//以下为SchemaVersion=2时才有的字段
	MediaType string
	Config    struct {
		MediaType string
		Size      int
		Digest    string
	}
	Layers []struct {
		MediaType string
		Size      int
		Digest    string
		Urls      interface{}
	}

	//以下未SchemaVersion=1时才有的字段
	Name         string
	Tag          string
	Architecture string
	FsLayers     []struct {
		BlobSum string
	}

	History []struct {
		History
		V1Compatibility string //非常恶心的未结构化的JSON字符串
	}

	//以下未SchemaVersion=1时，且请求Signed Manifest才有的字段
	Signatures []struct {
		Header struct {
			Jwk struct {
				Crv string
				Kid string
				Kty string
				X   string
				Y   string
			}
			Alg string
		}
		Signatures string
		Protected  string
	}
}

//获取Schema V1格式的Manifest
func (cli *Client) GetImageManifestV1(repo, ref, token string) (ImageManifest, error) {
	var info ImageManifest

	err := cli.RequestWithHeader("/v2/"+repo+"/manifests/"+ref, token, http.Header{}, &info)
	if err != nil {
		return info, err
	}

	//需要手动JSON解码
	for i, history := range info.History {
		json.Unmarshal([]byte(history.V1Compatibility), &info.History[i].History)
	}

	return info, err
}

//获取Schema V2格式的Manifest
func (cli *Client) GetImageManifestV2(repo, ref, token string) (ImageManifest, error) {
	var info ImageManifest

	header := http.Header{}
	header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	err := cli.RequestWithHeader("/v2/"+repo+"/manifests/"+ref, token, header, &info)
	return info, err
}
