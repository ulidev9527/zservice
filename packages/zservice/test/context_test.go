package test

import (
	"testing"
	"zserviceapps/packages/zservice"
)

func Test_ContextToJsonString(t *testing.T) {
	ctx := zservice.NewContext("")
	str := ctx.ToContextString()
	t.Log(str)
}
