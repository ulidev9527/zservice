package zservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func HttpRequestSend(ctx *Context, in *http.Request) ([]byte, *Error) {
	b, _ := json.Marshal(&ctx.ContextS2S)
	in.Header.Set(S_S2S_CTX, string(b))
	res, e := (&http.Client{}).Do(in) // 发起请求
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
	ctx.LogInfo("[zserver SendRequest] SUCC", res.Request.URL, string(body))
	return body, nil
}

// 发送 post 请求
func HttpPost(ctx *Context, url string, params map[string]any, header map[string]string) (body []byte, e *Error) {
	var bodyReader io.Reader
	logStr := ""
	if params != nil {
		for k, v := range params {
			vStr := fmt.Sprint(v)
			logStr = fmt.Sprint(logStr, "&", k, "=", vStr)
		}
		sByte, _ := json.Marshal(params)
		bodyReader = strings.NewReader(string(sByte))
	}

	req, _ := http.NewRequest(http.MethodPost, url, bodyReader)

	for k, v := range header {
		req.Header.Set(k, v)
	}
	ctx.LogInfof("[zserver.rest Post] %v %v", url, logStr)
	return HttpRequestSend(ctx, req)
}

// 发送 json 请求
func HttpPostJson(ctx *Context, url string, params map[string]any, header map[string]string) (body []byte, e *Error) {
	var bodyReader io.Reader
	logStr := ""
	if params != nil {
		for k, v := range params {
			vStr := fmt.Sprint(v)
			logStr = fmt.Sprint(logStr, "&", k, "=", vStr)
		}
		sByte, _ := json.Marshal(params)
		bodyReader = strings.NewReader(string(sByte))
	}

	req, _ := http.NewRequest(http.MethodPost, url, bodyReader)

	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.Header.Set("content-type", "application/json")

	ctx.LogInfof("[zserver.rest Post] %v %v", url, logStr)
	return HttpRequestSend(ctx, req)
}

// 发送 HttpGet 请求
func HttpGet(ctx *Context, url string, params map[string]any, header map[string]string) ([]byte, *Error) {
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
	return HttpRequestSend(ctx, req)
}
