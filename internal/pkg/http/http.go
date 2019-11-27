package http

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	apiHost string
	client  *http.Client
}

// 一个http请求, 忽略证书
func (c *Client) Request(path string, method string, post []byte, header http.Header, rsp interface{}) (code int, result []byte, err error) {
	url := c.apiHost + path
	var response *http.Response

	// 忽略https证书验证
	var req *http.Request
	if method == "GET" || method == "DELETE" {
		req, _ = http.NewRequest(method, url, nil)
		if header != nil {
			req.Header = header
		}
	} else if method == "POST" || method == "PUT" {
		req, _ = http.NewRequest(method, url, bytes.NewReader(post))
		if header != nil {
			req.Header = header
		}
	} else {
		req, _ = http.NewRequest(method, url, nil)
		if header != nil {
			req.Header = header
		}
	}
	// 只支持json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err = c.client.Do(req)
	if err != nil {
		return
	}

	defer response.Body.Close()
	code = response.StatusCode
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	result = body

	if rsp != nil {
		err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, rsp)
		if err != nil {
			err = errors.New(fmt.Sprintf("Unmarshal err: %v", err))
			return
		}
	}
	return
}

func (c *Client) Get(path string, header http.Header, rsp interface{}) (code int, result []byte, err error) {
	return c.Request(path, "GET", nil, header, rsp)
}

func (c *Client) Post(path string, post []byte, header http.Header, rsp interface{}) (code int, result []byte, err error) {
	return c.Request(path, "POST", post, header, rsp)
}

func (c *Client) Put(path string, post []byte, header http.Header, rsp interface{}) (code int, result []byte, err error) {
	return c.Request(path, "PUT", post, header, rsp)
}

func (c *Client) Delete(path string, header http.Header, rsp interface{}) (code int, result []byte, err error) {
	return c.Request(path, "POST", nil, header, rsp)
}

func NewClient(apiHost string) (c *Client, err error) {
	if apiHost == "" {
		err = errors.New("apiHost can't be empty")
		return
	}
	transport := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: transport, Timeout: 5 * time.Second}

	c = &Client{
		apiHost: apiHost,
		client:  client,
	}
	return
}
