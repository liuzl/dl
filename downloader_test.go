package dl

import (
	"flag"
	"testing"
)

var xUrl = flag.String("url", "http://m.newsmth.net", "url to fetch")

func TestDownload(t *testing.T) {
	flag.Parse()
	requestInfo := &HttpRequest{
		Url:      *xUrl,
		Method:   "GET",
		UseProxy: false,
		Platform: "google",
	}

	responseInfo := Download(requestInfo)
	if responseInfo.Error != nil {
		t.Error(responseInfo.Error)
	}
	t.Log(responseInfo.Text)
	t.Log(responseInfo.RemoteAddr)
}

func TestHTTPProxy(t *testing.T) {
	flag.Parse()
	requestInfo := &HttpRequest{
		Url:      *xUrl,
		Method:   "GET",
		UseProxy: true,
		Proxy:    "http://114.141.166.242:80",
		Platform: "google",
	}

	responseInfo := Download(requestInfo)
	if responseInfo.Error != nil {
		t.Error(responseInfo.Error)
	}
	t.Log(responseInfo.Text)
	t.Log(responseInfo.RemoteAddr)
}

func TestHTTPSProxy(t *testing.T) {
	flag.Parse()
	xUrl := "https://stackoverflow.com/"
	requestInfo := &HttpRequest{
		Url:      xUrl,
		Method:   "GET",
		UseProxy: true,
		Proxy:    "https://171.39.102.61:8123",
		Platform: "google",
	}

	responseInfo := Download(requestInfo)
	if responseInfo.Error != nil {
		t.Error(responseInfo.Error)
	}
	t.Log(responseInfo.Text)
	t.Log(responseInfo.RemoteAddr)
}

func TestSocks5Proxy(t *testing.T) {
	flag.Parse()
	requestInfo := &HttpRequest{
		Url:      *xUrl,
		Method:   "GET",
		UseProxy: true,
		Proxy:    "socks5://61.135.155.82:1080",
		Platform: "google",
	}

	responseInfo := Download(requestInfo)
	if responseInfo.Error != nil {
		t.Error(responseInfo.Error)
	}
	t.Log(responseInfo.Text)
	t.Log(responseInfo.RemoteAddr)
}
