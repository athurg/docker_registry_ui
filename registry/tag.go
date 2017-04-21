package registry

type TagsListResponse struct {
	BaseResponse
	Name string
	Tags []string
}

//获取指定仓库的标签列表
func (cli *Client) GetTags(repo, token string) (error, TagsListResponse) {
	var info TagsListResponse
	err := cli.Request("/v2/"+repo+"/tags/list", token, &info)
	return err, info
}
