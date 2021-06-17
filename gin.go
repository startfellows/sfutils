package sfutils

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gin-contrib/pprof"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	systemAccounts = gin.Accounts{
		"sflabs_system_account": "7Dnq9cgStvj7M1v8SvOaZZ0O",
	}
)

func ginPanicHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case error:
				// ugly as fuck
				xMsg := x.Error()
				if strings.Count(xMsg, "127.0.0.1") == 2 {
					if strings.Contains(xMsg, "write: broken pipe") ||
						strings.Contains(xMsg, "write: connection reset by peer") {
						break
					}
				}
				logrus.Errorln("[PANIC]", fmt.Sprintf("%s: %v\n%s", c.FullPath(), x.Error(), string(debug.Stack())))
			default:
				logrus.Errorln("[PANIC]", fmt.Sprintf("%s: %v\n%s", c.FullPath(), x, string(debug.Stack())))
			}
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}()

	c.Next()
}

// NewGin returns new gin.Engine object, ready to use.
// Configuration can be done via next environment variables:
// GIN_MODE - debug/release/test
// GIN_TEMPLATES_MODE - debug/release (if the variable is not set, templates are not loaded)
// GIN_TEMPLATES_PATH - absolute/relative path to html templates
// ENABLE_PPROF - enable/disable golang profiler on path /system/pprof
// ENABLE_PROMETHEUS - enable/disable gin metrics for prometheus on path /system/prometheus
// ENABLE_OPENAPI - enable/disable openapi schema on path /system/openapi
func NewGin() *gin.Engine {
	engine := gin.Default()

	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = true

	ginTemplatesMode := os.Getenv("GIN_TEMPLATES_MODE")
	if len(ginTemplatesMode) > 0 {
		var debugTemplates = true
		if ginTemplatesMode == "release" {
			debugTemplates = false
		}

		engine.HTMLRender = NewTemplateRender(os.Getenv("GIN_TEMPLATES_PATH"), ".tmpl", debugTemplates)
	}

	engine.Use(ginPanicHandler)

	systemGroup := engine.Group("/system", gin.BasicAuth(systemAccounts))

	if enable, _ := strconv.ParseBool(os.Getenv("ENABLE_PPROF")); enable {
		pprof.RouteRegister(systemGroup, "pprof")
	}

	if enable, _ := strconv.ParseBool(os.Getenv("ENABLE_PROMETHEUS")); enable {
		p := NewPrometheus("gin")
		p.UseRouterGroup(engine, systemGroup)
	}

	if enable, _ := strconv.ParseBool(os.Getenv("ENABLE_OPENAPI")); enable {
		systemGroup.GET("openapi", newOpenAPIHandler())
	}

	return engine
}
