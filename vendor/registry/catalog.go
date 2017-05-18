package registry

type CatalogResponseInfo struct {
	BaseResponse
	Repositories []string
}

//获取仓库列表
func (cli *Client) GetCatalog(token string) (error, CatalogResponseInfo) {
	var info CatalogResponseInfo
	err := cli.GetRequest("/v2/_catalog", token, nil, &info)
	return err, info
}
