package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/exvimmer/lenslocked/models"
	"github.com/gorilla/csrf"
)

type UsersTemplates struct {
	New    Template
	SignIn Template
}

type User struct {
	Templates   UsersTemplates
	UserService *models.UserService
}

func (u *User) New(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Email     string
		CsrfField template.HTML
	}{
		Email:     r.FormValue("email"),
		CsrfField: csrf.TemplateField(r),
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

func (u *User) SignIn(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Email string
	}{
		Email: r.FormValue("email"),
	}
	u.Templates.SignIn.Execute(w, data)
}

func (u *User) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	data := struct {
		email    string
		password string
	}{
		email:    r.FormValue("email"),
		password: r.FormValue("password"),
	}
	user, err := u.UserService.Authenticate(data.email, data.password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to authenticate", http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:     "email",
		Value:    data.email,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "user authenticated: %+v", user)
}

func (u *User) CurrentUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("email")
	if err != nil /* || cookie == nil */ {
		fmt.Fprintf(w, "the email value cannot be found in cookie")
		return
	}
	fmt.Fprintf(w, "email cookie: %s\n", cookie.Value)
}
