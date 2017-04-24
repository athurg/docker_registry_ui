package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BaseResponser interface {
	Error() error
}

type BaseResponse struct {
	Errors []struct {
		Code    string
		Message string
		Detail  interface{}
	}
}

func (resp *BaseResponse) Error() error {
	if len(resp.Errors) == 0 {
		return nil
	}

	var str string
	for _, errInfo := range resp.Errors {
		str += fmt.Sprintf("%s: %s (%s) ", errInfo.Code, errInfo.Message, errInfo.Detail)
	}

	return fmt.Errorf("%s", str)
}

//基础的HTTP请求，自动将Token加入Header中
func (cli *Client) Request(path, token string, result interface{}) error {
	return cli.RequestWithHeader(path, token, http.Header{}, result)
}

func (cli *Client) RequestWithHeader(path, token string, header http.Header, result interface{}) error {
	req, err := http.NewRequest("GET", cli.Addr+path, nil)
	req.Header = header
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

	if r, ok := result.(BaseResponser); ok {
		return r.Error()
	}

	return nil
}
