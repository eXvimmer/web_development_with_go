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
		New   Template
		Edit  Template
		Index Template
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
		http.Error(
			w,
			"you're not authorized to edit this gallery",
			http.StatusForbidden,
		)
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

func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
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
		http.Error(
			w,
			"you're not authorized to edit this gallery",
			http.StatusForbidden,
		)
		return
	}
	gallery.Title = r.FormValue("title")
	if gallery.Title == "" {
		http.Error(w, "Title cannot be empty", http.StatusInternalServerError)
		return
	}
	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(
		w,
		r,
		fmt.Sprintf("/galleries/%d/edit", gallery.Id),
		http.StatusFound,
	)
}

func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		Id    int
		Title string
	}
	var data struct {
		Galleries []Gallery
	}
	user := context.User(r.Context())
	// if user == nil {
	// 	http.Redirect(w, r, "/signin", http.StatusFound)
	// 	return
	// }
	galleries, err := g.GalleryService.ByUserId(user.Id)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			Id:    gallery.Id,
			Title: gallery.Title,
		})
	}
	g.Templates.Index.Execute(w, r, data)
}
