package registry

type TagsListResponse struct {
	BaseResponse
	Name string
	Tags []string
}

//获取指定仓库的标签列表
func (cli *Client) GetTags(repo, token string) (error, TagsListResponse) {
	var info TagsListResponse
	err := cli.GetRequest("/v2/"+repo+"/tags/list", token, nil, &info)
	return err, info
}
