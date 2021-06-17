package sfutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var openAPIData []byte

func SetOpenAPIData(data []byte) {
	openAPIData = data
}

func newOpenAPIHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "application/x-yaml", openAPIData)
	}
}
