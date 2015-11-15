package templates

import (
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// TemplateHandler represents a single template
type TemplateHandler struct {
	once     sync.Once
	Filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("../src/updown-kangaroo/templates", t.Filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}

	t.templ.Execute(w, data)
}
