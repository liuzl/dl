package dl

import (
	"context"
	"net/http"
)

type HttpRequest struct {
	Url        string                          `json:"url"`
	Method     string                          `json:"method"`
	PostData   string                          `json:"post_data"`
	UseProxy   bool                            `json:"use_proxy"`
	Proxy      string                          `json:"proxy"`
	Timeout    int                             `json:"timeout"`
	MaxLen     int64                           `json:"max_len"`
	Platform   string                          `json:"platform"`
	Retry      int                             `json:"retry"`
	Header     map[string]string               `json:"header"`
	Username   string                          `json:"username"`
	Password   string                          `json:"password"`
	ValidFuncs []func(resp *HttpResponse) bool `json:"-"`
	Ctx        context.Context                 `json:"-"`
	Jar        http.CookieJar                  `json:"-"`
}

type HttpResponse struct {
	Url        string          `json:"url"`
	Text       string          `json:"text"`
	Content    []byte          `json:"content"`
	Encoding   string          `json:"encoding"`
	StatusCode int             `json:"status_code"`
	Proxy      string          `json:"proxy"`
	RemoteAddr string          `json:"remote_addr"`
	Error      error           `json:"error"`
	Ctx        context.Context `json:"-"`
	Jar        http.CookieJar  `json:"-"`
}
