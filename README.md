# SF labs utils
The main purpose of this repository is to collect code that can be reused in one or more applications in order to speed up their development.

### Gin
Entry point for gin configuration is NewGin function.

Configuration can be done via next environment variables:
* GIN_MODE - debug/release/test
* GIN_TEMPLATES_MODE - debug/release (if the variable is not set, templates are not loaded)
* GIN_TEMPLATES_PATH - absolute/relative path to html templates
* ENABLE_PPROF - enable/disable golang profiler on path /system/pprof
* ENABLE_PROMETHEUS - enable/disable gin metrics for prometheus on path /system/prometheus
* ENABLE_OPENAPI - enable/disable openapi schema on path /system/openapi

### Gin templates
A small object that makes it easier to work with html templates in gin.
You can do without it if you don't need to render any html at all.

### Gin OpenAPI
A small hack to add openapi specification file to /system/openapi.

How to use it:
1. Set environment variable ```ENABLE_OPENAPI=true```
2. Add the following lines to your project:
```golang
//go:embed api.yml
var openAPIData []byte
```
3. Point sfutils to this variable:
```golang
sfutils.SetOpenAPIData(openAPIData)
```

### Remark about go pprof with basic auth
```bash
go tool pprof 'user:password@127.0.0.1:8080/system/pprof/heap'
go tool pprof 'user:password@127.0.0.1:8080/system/pprof/profile?seconds=30'
```