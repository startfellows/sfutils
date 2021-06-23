package auth

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func (a *Auth) RequireBearerSession(sessionObject interface{}) gin.HandlerFunc {
	sessionObjectType := reflect.TypeOf(sessionObject)

	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token = strings.TrimLeft(token, "Bearer ")

		i := reflect.New(sessionObjectType).Interface()

		err := a.Get(token, &i)
		if err != nil {
			if err != ErrSessionNotFound {
				logrus.Errorln("[Auth]", err)
			}

			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(a.prefix, i)
		c.Set(a.prefix+"_token", token)

		c.Next()
	}
}
