package dl

import (
	"testing"
)

func TestDownload(t *testing.T) {
	requestInfo := &HttpRequest{
		Url:      "http://m.newsmth.net",
		Method:   "GET",
		UseProxy: false,
		Platform: "mobile",
	}

	responseInfo := Download(requestInfo)
	if responseInfo.Error != nil {
		t.Error(responseInfo.Error)
	}
	t.Log(responseInfo.Text)
	t.Log(responseInfo.RemoteAddr)
}
