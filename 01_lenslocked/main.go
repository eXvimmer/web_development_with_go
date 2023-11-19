package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/exvimmer/lenslocked/controllers"
	"github.com/exvimmer/lenslocked/migrations"
	"github.com/exvimmer/lenslocked/models"
	"github.com/exvimmer/lenslocked/templates"
	"github.com/exvimmer/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	db, err := models.OpenDB(models.DefaultPostgresConfig())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	r.Route("/static", func(r chi.Router) {
		r.Handle("/*", http.StripPrefix("/static", fs))
	})

	sessionService := &models.SessionService{DB: db}

	usersC := controllers.User{
		Templates: controllers.UsersTemplates{
			New: views.Must(views.ParseFS(templates.FS, "signup.tmpl.html",
				"tailwind.tmpl.html")),
			SignIn: views.Must(views.ParseFS(templates.FS, "signin.tmpl.html",
				"tailwind.tmpl.html")),
		},
		UserService:    &models.UserService{DB: db},
		SessionService: sessionService,
	}

	r.Get("/",
		controllers.StaticHandler(
			views.Must(views.ParseFS(templates.FS, "home.tmpl.html",
				"tailwind.tmpl.html")), nil))

	r.Get("/contact",
		controllers.StaticHandler(
			views.Must(views.ParseFS(templates.FS, "contact.tmpl.html",
				"tailwind.tmpl.html")), nil))

	r.Get("/faq",
		controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.tmpl.html",
			"tailwind.tmpl.html"))))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/users/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfKey := "YoonjinMalena1992202313MustafaXl"
	csrfMW := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false), // TODO: set to true before deploying
	)

	fmt.Println(" 🚀 server is running on port :3000 ✅")
	err = http.ListenAndServe(":3000", csrfMW(umw.SetUser(r)))
	log.Fatal(err)
}
