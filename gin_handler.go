package sfutils

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type GinHandlerWithError func(c *gin.Context) error

type GinErrorHandler struct {
	errorTypes []reflect.Type
}

func (g *GinErrorHandler) ourErrorType(i interface{}) bool {
	errorType := reflect.TypeOf(i)

	for _, t := range g.errorTypes {
		fmt.Println("errorType:", errorType, t)
		if errorType == t {
			return true
		}
	}

	return false
}

func (g *GinErrorHandler) RegisterErrorType(i interface{}) {
	errorType := reflect.TypeOf(i)
	g.errorTypes = append(g.errorTypes, errorType)
}

func (g *GinErrorHandler) S(f GinHandlerWithError) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := f(c); err != nil {
			if g.ourErrorType(err) {
				c.JSON(http.StatusOK, gin.H{"error": err})
			} else {
				logrus.Errorf("WWW [%s] error: %s\n", c.FullPath(), err.Error())

				c.JSON(http.StatusInternalServerError, gin.H{})
			}
		}
	}
}

func NewGinErrorHandler() *GinErrorHandler {
	return &GinErrorHandler{}
}
