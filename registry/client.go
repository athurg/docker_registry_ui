package registry

type Client struct {
	Addr string
}

//Regisry客户端
func New(addr string) *Client {
	return &Client{Addr: addr}
}
