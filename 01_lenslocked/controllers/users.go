package controllers

import (
	"net/http"
)

type UsersTemplates struct {
	New Template
}

type Users struct {
	Templates UsersTemplates
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.Templates.New.Execute(w, nil)
}
