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
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err.Error())
		// TODO: send the right status code to the user
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	// TODO: return a session token
	fmt.Fprintf(w, "user created: %+v", user)
}
