package controllers

import (
	"net/http"

	"github.com/exvimmer/lenslocked/models"
)

type Galleries struct {
	Templates struct {
		New Template
	}
	GalleryService *models.GalleryService
}

func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: r.FormValue("title"),
	}
	g.Templates.New.Execute(w, r, data)
}
