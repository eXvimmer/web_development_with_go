package controllers

import (
	"net/http"

	"github.com/exvimmer/lenslocked/views"
)

func StaticHandler(t *views.Template, data any) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		t.Execute(w, data)
	}
}
