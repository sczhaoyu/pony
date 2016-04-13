package util

import (
	"io/ioutil"
	"net/http"
	u "net/url"
	"strings"
)

func HttpRequest(url, method string, prm map[string]string, header map[string]string) ([]byte, error) {
	uv := u.Values{}
	for k, v := range prm {
		uv.Set(k, v)
	}
	body := ioutil.NopCloser(strings.NewReader(uv.Encode()))
	client := &http.Client{}
	req, err := http.NewRequest(strings.ToUpper(method), url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req) //发送
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	data, derr := ioutil.ReadAll(resp.Body)
	if derr != nil {
		return nil, err
	}
	return data, derr
}
