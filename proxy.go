package dl

import (
	"errors"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/juju/errors"
)

const (
	proxyHost   = "http://127.0.0.1:8118"
	proxyGet    = "/get"
	proxyReport = "/bad"
)

type ProxyType int

const (
	Invalid ProxyType = iota
	HTTP
	HTTPS
	Socks5
)

func GetProxyType(proxy string) ProxyType {
	if strings.HasPrefix(proxy, "https") {
		return HTTPS
	}
	if strings.HasPrefix(proxy, "http") {
		return HTTP
	}
	if strings.HasPrefix(proxy, "socks5") {
		return Socks5
	}
	return Invalid
}

func GetProxy() (string, error) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	url := proxyHost + proxyGet
	for i := 0; i < 3; i++ {
		resp, err := client.Get(url)
		if err != nil {
			glog.Error("failed to get proxy, retry: ", i, "msg: ", err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			glog.Error("failed to read proxy, retry: ", i, "msg: ", err)
			continue
		}
		return "http://" + string(body), nil
	}
	return "", errors.Trace(errors.New("failed to get proxy finally"))
}

func ReportProxyStatus(proxy string) error {
	if len(proxy) == 0 {
		return nil
	}
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	reportUrl := proxyHost + proxyReport
	data := url.Values{}
	data.Add("p", proxy)
	resp, err := client.Post(reportUrl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return errors.Trace(err)
	}
	defer resp.Body.Close()
	return nil
}
