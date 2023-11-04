package controllers

import (
	"net/http"

	"github.com/exvimmer/lenslocked/views"
)

type UsersTemplates struct {
	New views.Template
}

type Users struct {
	Templates UsersTemplates
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.Templates.New.Execute(w, nil)
}
