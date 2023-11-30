package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/exvimmer/lenslocked/controllers"
	"github.com/exvimmer/lenslocked/migrations"
	"github.com/exvimmer/lenslocked/models"
	"github.com/exvimmer/lenslocked/static"
	"github.com/exvimmer/lenslocked/templates"
	"github.com/exvimmer/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

type CSRF struct {
	Key    string
	Secure bool
}

type config struct {
	Psql   *models.PostgressConfig
	Smtp   *models.SmtpConfig
	Csrf   CSRF
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	cfg := config{
		Psql: models.DefaultPostgresConfig(), // use this, if values are not set
		Smtp: &models.SmtpConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
		Csrf: CSRF{
			Key:    "YoonjinMalena1992202313MustafaXl",
			Secure: false, // TODO: set to true before deploying
		},
		Server: struct{ Address string }{
			Address: ":3000",
		},
	}
	// TODO: read psql, csrf, and server values from .env
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return cfg, err
	}
	cfg.Smtp.Port = port
	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	db, err := models.OpenDB(cfg.Psql)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	sessionService := &models.SessionService{DB: db}
	galleryService := &models.GalleryService{DB: db}

	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMW := csrf.Protect(
		[]byte(cfg.Csrf.Key),
		csrf.Secure(cfg.Csrf.Secure),
		csrf.Path("/"),
	)

	usersC := controllers.Users{
		Templates: controllers.UsersTemplates{
			New: views.Must(
				views.ParseFS(templates.FS, "signup.tmpl.html", "tailwind.tmpl.html"),
			),
			SignIn: views.Must(
				views.ParseFS(templates.FS, "signin.tmpl.html", "tailwind.tmpl.html"),
			),
			ForgotPassword: views.Must(
				views.ParseFS(
					templates.FS,
					"forgot_pw.tmpl.html",
					"tailwind.tmpl.html",
				),
			),
			CheckYourEmail: views.Must(
				views.ParseFS(
					templates.FS,
					"check_your_email.tmpl.html",
					"tailwind.tmpl.html",
				),
			),
			ResetPassword: views.Must(
				views.ParseFS(
					templates.FS,
					"reset_pw.tmpl.html",
					"tailwind.tmpl.html",
				),
			),
		},
		UserService:          &models.UserService{DB: db},
		SessionService:       sessionService,
		PasswordResetService: &models.PasswordResetService{DB: db},
		EmailService:         models.NewEmailService(cfg.Smtp),
	}
	galleryC := controllers.Galleries{
		GalleryService: galleryService,
	}
	galleryC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"galleries/new.tmpl.html",
		"tailwind.tmpl.html",
	))

	r := chi.NewRouter()
	r.Use(csrfMW, umw.SetUser)
	staticFileServer := http.FileServer(http.FS(static.FS))
	r.Handle("/static/*", http.StripPrefix("/static", staticFileServer))
	r.Get(
		"/",
		controllers.StaticHandler(
			views.Must(
				views.ParseFS(templates.FS, "home.tmpl.html", "tailwind.tmpl.html"),
			),
			nil,
		),
	)
	r.Get(
		"/contact",
		controllers.StaticHandler(
			views.Must(
				views.ParseFS(templates.FS, "contact.tmpl.html", "tailwind.tmpl.html"),
			),
			nil,
		),
	)
	r.Get(
		"/faq",
		controllers.FAQ(
			views.Must(
				views.ParseFS(templates.FS, "faq.tmpl.html", "tailwind.tmpl.html"),
			),
		),
	)
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.Route("/galleries", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/new", galleryC.New)
		})
	})
	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	fmt.Printf("🚀 server is running on port %s ✅\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	log.Fatal(err)
}
