package cns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// API响应的接口
type Responser interface {
	Error() error
}

// 基础的API响应结构
type BaseResponse struct {
	Code     int
	Message  string
	CodeDesc string
}

// API响应中识别并构造API错误的方法
func (resp BaseResponse) Error() error {
	if resp.Code == 0 {
		return nil
	}

	return fmt.Errorf("[%d](%s): %s", resp.Code, resp.CodeDesc, resp.Message)
}

// 云解析API请求的资源地址
var (
	_host     = "cns.api.qcloud.com"
	_SignHost = _host
	_uri      = _host + "/v2/index.php"
	_SignUri  = _uri
	_scheme   = "https"
)

func GetHost() string {
	return _host
}

// SetHost 设置API请求的域名
func SetHost(str string) {
	_host = str
	_uri = _host + "/v2/index.php"
}

// SetSignHost 设置API的签名HOST
func SetSignHost(host string) {
	_SignHost = host
	_SignUri = _SignHost + "/v2/index.php"
}

// SetHost 设置API请求的地址
func SetUri(str string) {
	_uri = str
}

// SetScheme 设置API请求的协议
func SetScheme(str string) {
	_scheme = str
}

func buildUri() string {
	return _scheme + "://" + _uri
}

// GET类型的API请求封装
func (cli *Client) requestGET(action string, param url.Values, respInfo interface{}) error {
	return cli.request("GET", action, param, nil, respInfo)
}

// API请求的封装（内建公共参数、签名的设置）
func (cli *Client) request(method, action string, param url.Values, body io.Reader, respInfo interface{}) error {
	if param == nil {
		param = url.Values{}
	}

	//设置公共参数
	param.Set("Action", action)
	param.Set("Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	param.Set("Nonce", "123456")
	param.Set("SecretId", cli.SecretId)

	sig := Signature(param, method, _SignUri, cli.SecretKey)
	param.Set("Signature", sig)

	req, err := http.NewRequest(method, buildUri()+"?"+param.Encode(), body)
	if err != nil {
		return fmt.Errorf("构建请求错误: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("执行请求错误: %s", err)
	}
	defer resp.Body.Close()

	w, ok := respInfo.(*bytes.Buffer)
	if ok {
		io.Copy(w, resp.Body)
		return nil
	}

	info, ok := respInfo.(Responser)
	if !ok {
		return fmt.Errorf("不可识别的响应结构参数")
	}

	err = json.NewDecoder(resp.Body).Decode(info)
	if err != nil {
		return fmt.Errorf("读取响应错误: %s", err)
	}

	return info.Error()
}
