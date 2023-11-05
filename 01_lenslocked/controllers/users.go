package controllers

import (
	"fmt"
	"net/http"
)

type UsersTemplates struct {
	New Template
}

type Users struct {
	Templates UsersTemplates
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Email string
	}{
		Email: r.FormValue("email"),
	}
	u.Templates.New.Execute(w, data)
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	// TODO: create the user
	fmt.Fprintf(w, "email: %s\npassword: %s", email, password)
}
