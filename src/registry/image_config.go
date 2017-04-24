package registry

import (
	"time"
)

type ContainerConfig struct {
	Hostname   string
	Domainname string
	Image      string
	User       string
	WorkingDir string

	AttachStdin  bool
	AttachStdout bool
	AttachStderr bool
	Tty          bool
	OpenStdin    bool
	StdinOnce    bool
	ArgsEscaped  bool

	Env        []string
	Cmd        []string
	OnBuild    []string
	Volumes    []string
	Entrypoint []string

	Labels map[string]string
}

//镜像配置对象
//    MIME为：application/vnd.docker.container.image.v1+json
//    定义参考: https://docs.docker.com/registry/spec/manifest-v2-2/
type ImageConfig struct {
	Os            string
	Architecture  string
	Author        string
	Container     string
	Created       time.Time
	DockerVersion string `json:"docker_version"`

	Config          ContainerConfig
	ContainerConfig ContainerConfig `json:"container_config"`

	History []struct {
		Created    time.Time
		CreatedBy  string `json:"created_by"`
		Author     string
		EmptyLayer bool `json:"empty_layer"`
	}

	RootFs struct {
		Type    string
		DiffIds []string `json:"diff_ids"`
	}
}

func (cli *Client) ImageConfigByDigest(repo, digest, token string) (ImageConfig, error) {
	var info ImageConfig

	err := cli.Request("/v2/"+repo+"/blobs/"+digest, token, &info)

	return info, err
}
