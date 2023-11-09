package controllers

import (
	"fmt"
	"net/http"

	"github.com/exvimmer/lenslocked/models"
)

type UsersTemplates struct {
	New Template
}

type User struct {
	Templates   UsersTemplates
	UserService *models.UserService
}

func (u *User) New(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Email string
	}{
		Email: r.FormValue("email"),
	}
	u.Templates.New.Execute(w, data)
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	// TODO: create the user
	fmt.Fprintf(w, "email: %s\npassword: %s", email, password)
}
