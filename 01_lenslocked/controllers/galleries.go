package controllers

import (
	"fmt"
	"net/http"

	"github.com/exvimmer/lenslocked/context"
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

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	data := struct {
		UserId int
		Title  string
	}{
		UserId: context.User(r.Context()).Id,
		Title:  r.FormValue("title"),
	}
	gallery, err := g.GalleryService.Create(data.Title, data.UserId)
	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.Id)
	http.Redirect(w, r, editPath, http.StatusFound)
}
