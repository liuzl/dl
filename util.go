package dl

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"github.com/saintfish/chardet"
	"math/rand"
)

var userAgents = map[string][]string{
	"pc": []string{
		"5.0 (X11; Linux x86_64; rv:29.0) Gecko/20100101 Firefox/29.0",
		"5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.137 Safari/537.36",
		"5.0 (Macintosh; Intel Mac OS X 10.9; rv:29.0) Gecko/20100101 Firefox/29.0",
		"5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.137 Safari/537.36",
		"Mac / Safari 7: 5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.75.14 (KHTML, like Gecko) Version/7.0.3 Safari/537.75.14",
		"5.0 (Windows NT 6.1; WOW64; rv:29.0) Gecko/20100101 Firefox/29.0",
		"5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.137 Safari/537.36",
		"4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)",
		"4.0 (compatible; MSIE 7.0; Windows NT 5.1)",
		"4.0 (compatible; MSIE 8.0; Windows NT 6.1; WOW64; Trident/4.0)",
		"5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)",
		"5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0)",
		"5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko"},
	"mobile": []string{
		"5.0 (Android; Mobile; rv:29.0) Gecko/29.0 Firefox/29.0",
		"5.0 (Linux; Android 4.4.2; Nexus 4 Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.114 Mobile Safari/537.36",
		"5.0 (iPad; CPU OS 7_0_4 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) CriOS/34.0.1847.18 Mobile/11B554a Safari/9537.53",
		"5.0 (iPad; CPU OS 7_0_4 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11B554a Safari/9537.53"}}

func GetUserAgent(platform string) string {
	if platform != "pc" && platform != "mobile" {
		platform = "pc"
	}
	return "Mozilla/" + userAgents[platform][rand.Intn(len(userAgents))]
}

func GetHeaders(platform string) map[string]string {
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language": "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3",
		"Accept-Encoding": "gzip, deflate",
		"User-Agent":      GetUserAgent(platform),
	}
	return headers
}

func GuessEncoding(content []byte) (string, error) {
	detector := chardet.NewHtmlDetector()
	res, err := detector.DetectBest(content)
	if err != nil {
		return "", err
	}
	return res.Charset, nil
}

func CompressPage(page []byte) string {
	var cPage bytes.Buffer
	cWriter := zlib.NewWriter(&cPage)
	cWriter.Write(page)
	cWriter.Close()
	b64Page := base64.StdEncoding.EncodeToString(cPage.Bytes())
	return b64Page
}
