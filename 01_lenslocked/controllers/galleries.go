package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/exvimmer/lenslocked/context"
	"github.com/exvimmer/lenslocked/errors"
	"github.com/exvimmer/lenslocked/models"
	"github.com/go-chi/chi/v5"
)

type Galleries struct {
	Templates struct {
		New  Template
		Edit Template
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

func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusNotFound)
		return
	}
	gallery, err := g.GalleryService.ById(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "gallery not found", http.StatusNotFound)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	user := context.User(r.Context())
	if user.Id != gallery.UserId {
		http.Error(w, "you're not authorized to edit this gallery", http.StatusForbidden)
		return
	}
	data := struct {
		Id    int
		Title string
	}{
		Id:    gallery.Id,
		Title: gallery.Title,
	}
	g.Templates.Edit.Execute(w, r, data)
}
