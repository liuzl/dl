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
