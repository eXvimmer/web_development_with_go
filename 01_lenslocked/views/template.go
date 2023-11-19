package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	myCtx "github.com/exvimmer/lenslocked/context"
	"github.com/exvimmer/lenslocked/models"
	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTmpl *template.Template
}

func ParseFS(fs fs.FS, patterns ...string) (*Template, error) {
	t := template.New(patterns[0]).Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField is not implemented")
			},
			"currentUser": func() (*models.User, error) {
				return nil, fmt.Errorf("currentUser is not implemented")
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
	// overwrite original ones
	tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *models.User {
			return myCtx.User(r.Context())
		},
	})
	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data) // to catch runtime errors
	if err != nil {
		log.Println("error executing template:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}
