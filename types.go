package dl

type HttpRequest struct {
	Url      string `json:"url"`
	Method   string `json:"method"`
	PostData string `json:"post_data"`
	UseProxy bool   `json:"use_proxy"`
	Proxy    string `json:"proxy"`
	Timeout  int    `json:"timeout"`
	MaxLen   int64  `json:"max_len"`
	Platform string `json:"platform"`
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
