package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTmpl *template.Template
}

func ParseFS(fs fs.FS, patterns ...string) (*Template, error) {
	t := template.New(patterns[0]).Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return `<!-- TODO: implement csrfField -->`
				// NOTE: this should be defined in the Execute method 
			},
		},
	)
	t, err := t.ParseFS(fs, patterns...)
	if err != nil {
		return &Template{}, fmt.Errorf("parsing FS: %w", err)
	}
	return &Template{
		htmlTmpl: t,
	}, nil
}

func (t *Template) Execute(w http.ResponseWriter, r *http.Request, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// NOTE: copy the template for each individual request to avoid race
	// condition. we could've used locks, but this is simpler.
	tpl, err := t.htmlTmpl.Clone()
	if err != nil {
		log.Printf("cloning template: %+v", err)
		http.Error(w, "cannot render the page", http.StatusInternalServerError)
		return
	}
	tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
	})
	err = tpl.Execute(w, data)
	if err != nil {
		log.Println("error executing template:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}
