package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	myCtx "github.com/exvimmer/lenslocked/context"
	"github.com/exvimmer/lenslocked/models"
)

type UsersTemplates struct {
	New            Template
	SignIn         Template
	ForgotPassword Template
	CheckYourEmail Template
	ResetPassword  Template
}

type Users struct {
	Templates            UsersTemplates
	UserService          *models.UserService
	SessionService       *models.SessionService
	EmailService         *models.EmailService
	PasswordResetService *models.PasswordResetService
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Email string
	}{
		Email: r.FormValue("email"),
	}
	u.Templates.New.Execute(w, r, data)
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
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

func (u *Users) SignIn(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Email string
	}{
		Email: r.FormValue("email"),
	}
	u.Templates.SignIn.Execute(w, r, data)
}

func (u *Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
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

func (u *Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := myCtx.User(r.Context())
	// NOTE: checking user == nil is already done in RequireUser middleware, so
	// we don't need to check it again.
	fmt.Fprintf(w, "current user: %s\n", user.Email)
}

func (u *Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
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

func (u *Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u *Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		// TODO: handle other cases, like non-existing email address
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetUrl := "https://lenslocked.com/reset-pw?" + vals.Encode()
	err = u.EmailService.ForgotPassword(data.Email, resetUrl)
	if err != nil {
		// TODO: handle other cases, like non-existing email address
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	// WARN: don't render the reset token in the template. users should be able
	// to confirm that they have access to the email account to verify their
	// identity.
	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (u *Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u *Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")
	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	err = u.UserService.UpdatePassword(user.Id, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.Id)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
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
