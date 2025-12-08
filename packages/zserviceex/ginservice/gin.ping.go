package ginservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func gin_ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
