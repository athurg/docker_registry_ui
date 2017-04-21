package registry

type CatalogResponseInfo struct {
	Repositories []string
}

//获取仓库列表
func (cli *Client) GetCatalog(token string) (error, CatalogResponseInfo) {
	var info CatalogResponseInfo
	err := cli.Request("/v2/_catalog", token, &info)
	return err, info
}
