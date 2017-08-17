package dl

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Download(requestInfo *HttpRequest) *HttpResponse {
	var timeout time.Duration
	if requestInfo.Timeout > 0 {
		timeout = time.Duration(requestInfo.Timeout) * time.Second
	} else {
		timeout = 30 * time.Second
	}
	client := &http.Client{
		Timeout: timeout,
	}
	responseInfo := &HttpResponse{
		Url: requestInfo.Url,
	}
	transport := http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	}

	//proxy
	if requestInfo.UseProxy {
		var proxy string
		var err error
		if len(requestInfo.Proxy) > 0 {
			proxy = requestInfo.Proxy
		} else {
			proxy, err = GetProxy()
			if err != nil {
				responseInfo.Error = err
				return responseInfo
			}
		}
		responseInfo.Proxy = proxy
		urlProxy, err := url.Parse(proxy)
		if err != nil {
			responseInfo.Error = errors.New(fmt.Sprintf("failed to parse proxy: %s", proxy))
			return responseInfo
		}
		transport.Proxy = http.ProxyURL(urlProxy)
	}

	client.Transport = &transport

	req, err := http.NewRequest(requestInfo.Method, requestInfo.Url, strings.NewReader(requestInfo.PostData))
	if err != nil {
		responseInfo.Error = err
		return responseInfo
	}
	headers := GetHeaders(requestInfo.Platform)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if requestInfo.Method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			responseInfo.RemoteAddr = connInfo.Conn.RemoteAddr().String()
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		responseInfo.Error = err
		return responseInfo
	}

	responseInfo.StatusCode = resp.StatusCode
	defer resp.Body.Close()

	var contentLen int64
	contentLen, err = strconv.ParseInt(resp.Header.Get("content-length"), 10, 64)
	if err != nil {
		//
	} else if requestInfo.MaxLen > 0 && contentLen > requestInfo.MaxLen {
		responseInfo.Error = errors.New("reponse size too large")
		return responseInfo
	}

	var readLen int64 = 0
	respBuf := bytes.NewBuffer([]byte{})
	for {
		readData := make([]byte, 4096)
		length, err := resp.Body.Read(readData)
		respBuf.Write(readData[:length])
		readLen += int64(length)
		if err != nil {
			if err == io.EOF {
				break
			}
			responseInfo.Error = errors.New("reponse size too large - count")
			return responseInfo
		}
	}
	content := respBuf.Bytes()

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		var reader io.ReadCloser
		contentReader := bytes.NewReader(content)
		if reader, err = gzip.NewReader(contentReader); err != nil {
			responseInfo.Error = err
			return responseInfo
		}
		defer reader.Close()
		responseInfo.Content, _ = ioutil.ReadAll(reader)
	case "deflate":
		var reader io.ReadCloser
		contentReader := bytes.NewReader(content)
		if reader, err = zlib.NewReader(contentReader); err != nil {
			if err == zlib.ErrHeader {
				// raw defalte, no zlib header
				reader = flate.NewReader(bytes.NewReader(content))
			} else {
				responseInfo.Error = err
				return responseInfo
			}
		}
		defer reader.Close()
		responseInfo.Content, _ = ioutil.ReadAll(reader)
	default:
		responseInfo.Content = content
	}

	var encoding string
	encoding, err = GuessEncoding(responseInfo.Content)
	if err != nil {
		//
		responseInfo.Text = string(responseInfo.Content)
		responseInfo.Encoding = ""
		return responseInfo
	}
	encoder := mahonia.NewDecoder(encoding)
	if encoder == nil {
		responseInfo.Text = string(responseInfo.Content)
		responseInfo.Encoding = ""
		return responseInfo
	}
	responseInfo.Text = encoder.ConvertString(string(responseInfo.Content))
	responseInfo.Encoding = encoding
	return responseInfo
}
