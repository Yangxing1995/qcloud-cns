package cns

// 云解析API请求的客户端
type Client struct {
	SecretId  string
	SecretKey string
	Host      string
	SignHost  string
	Scheme    string
}

// 云解析API请求的资源地址
const (
	_defaultHost   = "cns.api.qcloud.com"
	_defaultScheme = "https"
	_indexURL      = "/v2/index.php"
)

func New(secretId, secretKey string) *Client {
	return &Client{
		SecretId:  secretId,
		SecretKey: secretKey,
		Host:      _defaultHost,
		SignHost:  _defaultHost,
		Scheme:    _defaultScheme,
	}
}

func (cli *Client) SetHost(host string) *Client {
	cli.Host = host
	return cli
}

func (cli *Client) SetScheme(str string) *Client {
	cli.Scheme = str
	return cli
}

func (cli *Client) SetSignHost(host string) *Client {
	cli.SignHost = host
	return cli
}

func (cli *Client) buildUri() string {
	return cli.Scheme + "://" + cli.Host + _indexURL
}

func (cli *Client) buildSignUri() string {
	return cli.SignHost + _indexURL
}
