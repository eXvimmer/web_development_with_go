package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/exvimmer/lenslocked/context"
	"github.com/exvimmer/lenslocked/errors"
	"github.com/exvimmer/lenslocked/models"
	"github.com/go-chi/chi/v5"
)

type Image struct {
	GalleryId       int
	Filename        string
	FilenameEscaped string
}

type Galleries struct {
	Templates struct {
		New   Template
		Edit  Template
		Index Template
		Show  Template
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

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(w, r)
	if err != nil {
		return // g.galleryById handles the rendering
	}
	data := struct {
		Id     int
		Title  string
		Images []Image
	}{
		Id:    gallery.Id,
		Title: gallery.Title,
	}
	images, err := g.GalleryService.Images(gallery.Id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryId:       image.GalleryId,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	g.Templates.Show.Execute(w, r, data)
}

func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(w, r, userMustOwnGallery)
	if err != nil {
		return // g.galleryById handles the rendering
	}
	data := struct {
		Id     int
		Title  string
		Images []Image
	}{
		Id:    gallery.Id,
		Title: gallery.Title,
	}
	images, err := g.GalleryService.Images(gallery.Id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryId:       image.GalleryId,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	g.Templates.Edit.Execute(w, r, data)
}

func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(w, r, userMustOwnGallery)
	if err != nil {
		return // g.galleryById handles the rendering
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

func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(w, r, userMustOwnGallery)
	if err != nil {
		return // g.galleryById handles the rendering
	}
	err = g.GalleryService.Delete(gallery.Id)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g *Galleries) Image(w http.ResponseWriter, r *http.Request) {
	filename := g.filename(r)
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid gallery id", http.StatusNotFound)
		return
	}
	image, err := g.GalleryService.Image(id, filename)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, image.Path)
}

func (g *Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	filename := g.filename(r)
	err = g.GalleryService.DeleteImage(gallery.Id, filename)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(
		w,
		r,
		fmt.Sprintf("/galleries/%d/edit", gallery.Id),
		http.StatusFound,
	)
}

// NOTE: read Rob Pike's blog post about self referential functions and the
// design of options.
// https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
type galleryOpt func(http.ResponseWriter, *http.Request, *models.Gallery) error

func (g *Galleries) galleryById(
	w http.ResponseWriter,
	r *http.Request,
	opts ...galleryOpt,
) (*models.Gallery, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.GalleryService.ById(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "gallery not found", http.StatusNotFound)
			return nil, err
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return nil, err
	}
	for _, opt := range opts {
		err := opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}
	return gallery, nil
}

func (g *Galleries) filename(r *http.Request) string {
	return filepath.Base(chi.URLParam(r, "filename"))
}

func userMustOwnGallery(
	w http.ResponseWriter,
	r *http.Request,
	gallery *models.Gallery,
) error {
	user := context.User(r.Context())
	if user.Id != gallery.UserId {
		http.Error(
			w,
			"you're not authorized to edit this gallery",
			http.StatusForbidden,
		)
		return fmt.Errorf("not authorized")
	}
	return nil
}
