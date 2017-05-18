package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	Addr string
}

//Regisry客户端
func New(addr string) *Client {
	return &Client{Addr: addr}
}

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

func (cli *Client) requestWithToken(method, path, token string, header http.Header) (*http.Response, error) {
	req, err := http.NewRequest(method, cli.Addr+path, nil)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	} else {
		req.Header = http.Header{}
	}

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}

	//4xx错误会携带错误信息
	if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
		var errorInfo BaseResponse
		err = json.NewDecoder(resp.Body).Decode(&errorInfo)
		resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("读取错误:", err)
		}

		return nil, errorInfo.Error()
	}

	return resp, err
}

//基础的HTTP请求，自动将Token加入Header中
func (cli *Client) GetRequest(path, token string, header http.Header, result interface{}) error {
	resp, err := cli.requestWithToken("GET", path, token, header)
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
