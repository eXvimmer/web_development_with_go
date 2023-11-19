package controllers

import (
	"fmt"
	"net/http"

	myCtx "github.com/exvimmer/lenslocked/context"
	"github.com/exvimmer/lenslocked/models"
)

type UsersTemplates struct {
	New    Template
	SignIn Template
}

type User struct {
	Templates      UsersTemplates
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u *User) New(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Email string
	}{
		Email: r.FormValue("email"),
	}
	u.Templates.New.Execute(w, r, data)
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.Id)
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u *User) SignIn(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Email string
	}{
		Email: r.FormValue("email"),
	}
	u.Templates.SignIn.Execute(w, r, data)
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
	session, err := u.SessionService.Create(user.Id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u *User) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := myCtx.User(r.Context())
	// NOTE: checking user == nil is already done in RequireUser middleware, so
	// we don't need to check it again.
	fmt.Fprintf(w, "current user: %s\n", user.Email)
}

func (u *User) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := getCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw *UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := getCookie(r, CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := umw.SessionService.User(cookie)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		newCtx := myCtx.WithUser(r.Context(), user)
		r = r.WithContext(newCtx)
		next.ServeHTTP(w, r)
	})
}

// this method assumes to be called after SetUser method
func (umw *UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := myCtx.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
