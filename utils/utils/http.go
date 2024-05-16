package utils

import (
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
)

// http工具
type httpUtils struct {
}

func (u *utils) NewHttpUtils() *httpUtils {
	return &httpUtils{}
}

func (h *httpUtils) Client() (resp *resty.Request) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli := &http.Client{
		Transport: tr,
	}
	resp = resty.NewWithClient(cli).R()
	return
}

func (h *httpUtils) ProxyClient(proxyUrl string) (resp *resty.Request) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	uri, err := url.Parse(proxyUrl)
	if err == nil {
		tr.Proxy = http.ProxyURL(uri)
	}
	cli := &http.Client{
		Transport: tr,
	}
	resp = resty.NewWithClient(cli).R()
	return
}
