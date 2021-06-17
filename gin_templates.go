package sfutils

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/render"
	"github.com/sirupsen/logrus"
)

type TemplateRender struct {
	templates      map[string]*template.Template
	templatesMutex sync.Mutex

	templatesDir string
	ext          string
	debug        bool
}

func (r *TemplateRender) Reload() {
	r.templatesMutex.Lock()
	defer r.templatesMutex.Unlock()

	r.templates = map[string]*template.Template{}
}

func (r *TemplateRender) GetTemplate(name string) *template.Template {
	r.templatesMutex.Lock()
	defer r.templatesMutex.Unlock()

	// Check if gin is running in debug mode and load the templates accordingly
	tpl := r.templates[name]
	if tpl == nil {
		tpl = r.loadTemplate(name)

		if !r.debug {
			r.templates[name] = tpl
		}
	}

	return tpl
}

func (r *TemplateRender) Instance(name string, data interface{}) render.Render {
	tpl := r.GetTemplate(name)

	return render.HTML{
		Template: tpl,
		Data:     data,
	}
}

func (r *TemplateRender) loadTemplate(name string) *template.Template {
	// get all templates from includes/
	var includes []string
	filepath.Walk(path.Join(r.templatesDir, "includes"),
		func(path string, f os.FileInfo, err error) error {
			if strings.HasSuffix(path, r.ext) {
				includes = append(includes, path)
			}

			return nil
		})

	// get template
	file := path.Join(r.templatesDir, name+r.ext)
	//tpl, _ := template.ParseFiles(append([]string{file}, r.Includes...)...)

	tpl, err := template.New(path.Base(file)).
		Funcs(ginFuncMap).
		ParseFiles(append([]string{file}, includes...)...)

	if err != nil {
		logrus.Errorln("[GIN Template]", fmt.Sprintf("%s: %s", name, err.Error()))
	}

	return tpl
}

func NewTemplateRender(templatesDir string, ext string, debug bool) *TemplateRender {
	r := &TemplateRender{
		templates: map[string]*template.Template{},

		templatesDir: templatesDir,
		ext:          ext,
		debug:        debug,
	}

	return r
}
