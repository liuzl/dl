package dl

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	Proxy "golang.org/x/net/proxy"
)

func DownloadUrl(url string) *HttpResponse {
	return Download(&HttpRequest{Url: url})
}

func DownloadUrlWithProxy(url string) *HttpResponse {
	return Download(&HttpRequest{Url: url, UseProxy: true, Retry: 5})
}

func Download(requestInfo *HttpRequest) *HttpResponse {
	if requestInfo.Retry == 0 {
		requestInfo.Retry = 1
	}
	var resp *HttpResponse
	for i := 0; i < requestInfo.Retry; i++ {
		resp = downloadOnce(requestInfo)
		if resp != nil && (resp.ctx.Err() == context.Canceled || resp.ctx.Err() == context.DeadlineExceeded) {
			return resp
		}
		if resp == nil || resp.Error != nil {
			time.Sleep(time.Second * time.Duration(rand.Intn(2)+1))
			continue
		}
		respValid := true
		if requestInfo.ValidFuncs != nil {
			for _, validFunc := range requestInfo.ValidFuncs {
				if !validFunc(resp) {
					respValid = false
					break
				}
			}
		}
		if !respValid {
			time.Sleep(time.Second * time.Duration(rand.Intn(2)+1))
			continue
		} else {
			break
		}
	}
	return resp
}
func downloadOnce(requestInfo *HttpRequest) *HttpResponse {
	var timeout time.Duration
	if requestInfo.Timeout > 0 {
		timeout = time.Duration(requestInfo.Timeout) * time.Second
	} else {
		timeout = 30 * time.Second
	}
	client := &http.Client{Timeout: timeout}
	responseInfo := &HttpResponse{Url: requestInfo.Url}
	transport := http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	}

	//proxy
	needReport := false
	var proxy string
	if requestInfo.UseProxy {
		var err error
		if len(requestInfo.Proxy) > 0 {
			proxy = requestInfo.Proxy
		} else {
			needReport = true
			proxy, err = GetProxy()
			if err != nil {
				responseInfo.Error = err
				return responseInfo
			}
		}
		responseInfo.Proxy = proxy
		urlProxy, err := url.Parse(proxy)
		if err != nil {
			responseInfo.Error = fmt.Errorf("failed to parse proxy: %s", proxy)
			return responseInfo
		}
		proxyType := GetProxyType(proxy)
		switch proxyType {
		case Invalid:
			responseInfo.Error = fmt.Errorf("invalid proxy type, proxy: %s", proxy)
			return responseInfo
		case Socks5:
			proxyAddr := strings.Trim(proxy, "socks5://")
			dialer, err := Proxy.SOCKS5("tcp", proxyAddr, nil, Proxy.Direct)
			if err != nil {
				responseInfo.Error = fmt.Errorf(
					"failed to set socks5 proxy, proxy: %s, msg: %s", proxy, err)
				return responseInfo
			}
			transport.Dial = dialer.Dial
		case HTTP, HTTPS:
			transport.Proxy = http.ProxyURL(urlProxy)
		}
	}

	client.Transport = &transport

	req, err := http.NewRequest(
		requestInfo.Method, requestInfo.Url, strings.NewReader(requestInfo.PostData))
	if err != nil {
		responseInfo.Error = err
		return responseInfo
	}
	if requestInfo.ctx != nil {
		req = req.WithContext(requestInfo.ctx)
	}
	headers := GetHeaders(requestInfo.Platform)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	for k, v := range requestInfo.Header {
		req.Header.Set(k, v)
	}
	if requestInfo.Method == "POST" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			responseInfo.RemoteAddr = connInfo.Conn.RemoteAddr().String()
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	responseInfo.ctx = req.Context()
	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		if needReport {
			ReportProxyStatus(proxy)
		}
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
		responseInfo.Error = fmt.Errorf("reponse size too large")
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
			responseInfo.Error = fmt.Errorf("reponse size too large - count")
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
	if encoding, err = GuessEncoding(responseInfo.Content); err != nil {
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
