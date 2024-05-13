package zservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func RequestSend(ctx *Context, req *http.Request) ([]byte, *Error) {
	b, _ := json.Marshal(&ctx.ContextTrace)
	req.Header.Set(S_TraceKey, string(b))
	res, e := (&http.Client{}).Do(req) // 发起请求
	if e != nil {
		return nil, NewError(e)
	}

	if res.StatusCode != http.StatusOK {
		return nil, NewError("REQ FAIL", res.Request.URL)
	}

	body, e := io.ReadAll(res.Body)
	if e != nil {
		return nil, NewError(e)
	}
	ctx.LogInfo("[zserver SendRequest] RES SUCC", string(body))
	return body, nil
}

// 发送 post 请求
func Post(ctx *Context, url string, params *map[string]any, header *map[string]string) (body []byte, e *Error) {
	var bodyReader io.Reader
	logStr := ""
	if params != nil {
		for k, v := range *params {
			vStr := fmt.Sprint(v)
			logStr = fmt.Sprint(logStr, "&", k, "=", vStr)
		}
		sByte, _ := json.Marshal(params)
		bodyReader = strings.NewReader(string(sByte))
	}

	req, _ := http.NewRequest(http.MethodPost, url, bodyReader)

	if header != nil {
		for k, v := range *header {
			req.Header.Set(k, v)
		}
	}
	ctx.LogInfof("[zserver.rest Post] %v %v", url, logStr)
	return RequestSend(ctx, req)
}

// 发送 json 请求
func PostJson(ctx *Context, url string, params *map[string]any, header *map[string]string) (body []byte, e *Error) {
	if header == nil {
		header = &map[string]string{}
	}

	var bodyReader io.Reader
	logStr := ""
	if params != nil {
		for k, v := range *params {
			vStr := fmt.Sprint(v)
			logStr = fmt.Sprint(logStr, "&", k, "=", vStr)
		}
		sByte, _ := json.Marshal(params)
		bodyReader = strings.NewReader(string(sByte))
	}

	req, _ := http.NewRequest(http.MethodPost, url, bodyReader)

	if header != nil {
		for k, v := range *header {
			req.Header.Set(k, v)
		}
	}

	req.Header.Set("content-type", "application/json")

	ctx.LogInfof("[zserver.rest Post] %v %v", url, logStr)
	return RequestSend(ctx, req)
}

// 发送 Get 请求
func Get(ctx *Context, url string, params *map[string]any, header *map[string]string) (body []byte, e *Error) {
	logStr := ""
	if params != nil {
		for k, v := range *params {
			logStr = fmt.Sprint(logStr, "&", k, "=", v)
		}
	}
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%v?%v", url, logStr[1:]), nil)

	if header != nil {
		for k, v := range *header {
			req.Header.Set(k, v)
		}
	}
	return RequestSend(ctx, req)
}
