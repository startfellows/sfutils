package sfutils

import "html/template"

var ginFuncMap = template.FuncMap{
	"noescape": templateNoEscape,
	"inc":      templateInc,
	"dec":      templateDec,
}

func templateNoEscape(str string) template.HTML {
	return template.HTML(str)
}

func templateInc(a, b int) int {
	return a + b
}

func templateDec(a, b int) int {
	return a - b
}
