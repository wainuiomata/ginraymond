// Package ginraymond is a custom template renderer that can be used with
// the Gin web framework: https://github.com/gin-gonic/gin, it adds support
// for Handlebars style templates to your Gin application. It uses the Raymond
// template library for this: https://github.com/aymerick/raymond, which
// implements Handlebars style templates in pure Go.
package ginraymond

import (
	"net/http"
	"path"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

// RaymondRender is a custom Gin template renderer using Raymond.
type RaymondRender struct {
	Options  *RenderOptions
	Cache    *TemplateCache
	Template *raymond.Template
	Context  interface{}
}

// New creates a new RaymondRender instance with custom Options.
func New(options *RenderOptions) *RaymondRender {
	return &RaymondRender{
		Options: options,
		Cache: &TemplateCache{
			TemplateDir: options.TemplateDir,
			cache:       make(map[string]*raymond.Template),
		},
	}
}

// Default creates a RaymondRender instance with default options.
func Default() *RaymondRender {
	return New(DefaultOptions())
}

// Instance should return a new RaymondRender struct for the current request
// and prepare the *raymond.Template that will be used for this response.
// Because of GIN's design, the Instance method cannot return an error type.
// This means that we have to panic on any errors which is not ideal in Go.
func (r RaymondRender) Instance(name string, data interface{}) render.Render {
	var template *raymond.Template
	filename := path.Join(r.Options.TemplateDir, name)

	// always read template files from disk if in debug mode, use cache otherwise.
	if gin.Mode() == "debug" {
		template = MustLoadTemplate(filename)
	} else {
		template = r.Cache.MustGet(filename)
	}

	return RaymondRender{
		Template: template,
		Context:  data,
		Options:  r.Options,
	}
}

// Render should write the content type, then render the template to the response.
func (r RaymondRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	output, err := r.Template.Exec(r.Context)
	w.Write([]byte(output))
	return err
}

// WriteContentType writes header information about content RaymondRender outputs.
// This will now implement gin's render.Render interface.
func (r RaymondRender) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = r.Options.ContentType
	}
}
