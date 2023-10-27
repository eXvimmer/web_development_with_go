package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmlTmpl *template.Template
}

func (t *Template) Execute(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmlTmpl.Execute(w, data)
	if err != nil {
		log.Println("error executing template:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func Parse(filepath string) (*Template, error) {
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %w", err)
	}
	return &Template{htmlTmpl: tmpl}, nil
}

func ParseFS(fs fs.FS, pattern string) (*Template, error) {
	t, err := template.ParseFS(fs, pattern)
	if err != nil {
		return &Template{}, fmt.Errorf("parsing FS: %w", err)
	}
	return &Template{
		htmlTmpl: t,
	}, nil
}

func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}
