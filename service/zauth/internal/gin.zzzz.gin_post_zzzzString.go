package internal

import (
	"net/http"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

func gin_post_zzzzString(c *gin.Context) {
	zctx := ginservice.GetCtxEX(c)

	file, e := c.FormFile("file")
	if e != nil {
		zctx.LogError(e)
		c.JSON(http.StatusOK, gin.H{"code": zglobal.Code_ErrorBreakoff})
		return
	}

	if e := c.SaveUploadedFile(file, FI_ZZZZStringFile); e != nil {
		zctx.LogError(e)
		c.JSON(http.StatusOK, gin.H{"code": zglobal.Code_ErrorBreakoff})
		return
	} else if e := ZZZZString.Reload(zctx); e != nil {
		zctx.LogError(e)
		c.JSON(http.StatusOK, gin.H{"code": e.GetCode()})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": zglobal.Code_SUCC})
	}
}
