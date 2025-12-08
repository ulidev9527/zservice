package zservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func HttpRequestSend_Old(ctx *Context, in *http.Request) ([]byte, *Error) {
	res, e := (&http.Client{}).Do(in)
	if e != nil {
		return nil, NewError(e)
	}

	if res.StatusCode != http.StatusOK {
		return nil, NewError("[zserver SendRequest] FAIL", res.Request.URL)
	}

	body, e := io.ReadAll(res.Body)
	if e != nil {
		return nil, NewError(e)
	}
	return body, nil
}

// 发送 post 请求
func HttpPost_Old(ctx *Context, url string, params map[string]any, header map[string]string) (body []byte, e *Error) {
	var bodyReader io.Reader
	if params != nil {
		sByte, _ := json.Marshal(params)
		bodyReader = strings.NewReader(string(sByte))
	}

	req, _ := http.NewRequest(http.MethodPost, url, bodyReader)

	for k, v := range header {
		req.Header.Set(k, v)
	}
	return HttpRequestSend_Old(ctx, req)
}

// 发送 HttpGet_Old 请求
func HttpGet_Old(ctx *Context, url string, params map[string]any, header map[string]string) ([]byte, *Error) {
	paramsStr := ""
	for k, v := range params {
		paramsStr = fmt.Sprint(paramsStr, "&", k, "=", v)
	}

	if len(paramsStr) > 0 {
		if !strings.Contains(url, "?") {
			url = url + "?"
		}
		url = url + paramsStr
	}

	req, e := http.NewRequest(http.MethodGet, url, nil)
	if e != nil {
		return nil, NewError(e)
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}
	return HttpRequestSend_Old(ctx, req)
}
