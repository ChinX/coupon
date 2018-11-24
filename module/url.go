package module

import "net/url"

var HostName = "https://piao.windup.cn"

func urlFormat(path string) string {
	rawURL, err := url.Parse(path)
	if err != nil || rawURL.Host == ""{
		rawURL, _ = url.Parse(HostName+path)
	}
	return rawURL.String()
}
