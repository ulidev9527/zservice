package zservice

import "resty.dev/v3"

var restyClient = resty.New()

func HttpPostPB(path string, body []byte) []byte {
	if res, e := restyClient.R().
		SetContentType(Http_ContentType_ProtoBuf).
		SetBody(body).
		Post(path); e != nil {
		LogError(path, e)
		return nil
	} else {
		return res.Bytes()
	}
}

// 带上下文
func HttpPostPBCtx(ctx *Context, path string, body []byte) []byte {
	if res, e := restyClient.R().
		SetContentType(Http_ContentType_ProtoBuf).
		SetBody(body).
		SetHeader("Zctx", ctx.ToContextString()).
		Post(path); e != nil {
		LogError(path, e)
		return nil
	} else {
		return res.Bytes()
	}
}
