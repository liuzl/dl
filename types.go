package dl

import (
	"context"
)

type HttpRequest struct {
	Url        string                          `json:"url"`
	Method     string                          `json:"method"`
	PostData   string                          `json:"post_data"`
	UseProxy   bool                            `json:"use_proxy"`
	Proxy      string                          `json:"proxy"`
	MaxLen     int64                           `json:"max_len"`
	Platform   string                          `json:"platform"`
	Retry      int                             `json:"retry"`
	Header     map[string]string               `json:"header"`
	ValidFuncs []func(resp *HttpResponse) bool `json:"-"`
	ctx        context.Context                 `json:"-"`
}

type HttpResponse struct {
	Url        string            `json:"url"`
	Text       string            `json:"text"`
	Content    []byte            `json:"content"`
	Encoding   string            `json:"encoding"`
	StatusCode int               `json:"status_code"`
	Proxy      string            `json:"proxy"`
	Cookies    map[string]string `json:"cookies"`
	RemoteAddr string            `json:"remote_addr"`
	Error      error             `json:"error"`
}
