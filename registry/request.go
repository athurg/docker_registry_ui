package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//基础的HTTP请求，自动将Token加入Header中
func (cli *Client) Request(path, token string, result interface{}) error {
	req, err := http.NewRequest("GET", cli.Addr+path, nil)
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求错误:", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("读取错误:", err)
	}

	return nil
}
