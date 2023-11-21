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

	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMW := csrf.Protect(
		[]byte(cfg.Csrf.Key),
		csrf.Secure(cfg.Csrf.Secure),
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
		},
		UserService:          &models.UserService{DB: db},
		SessionService:       sessionService,
		PasswordResetService: &models.PasswordResetService{DB: db},
		EmailService:         models.NewEmailService(cfg.Smtp),
	}

	fs := http.FileServer(http.Dir("./static"))
	r := chi.NewRouter()
	r.Use(csrfMW, umw.SetUser)
	r.Route(
		"/static",
		func(r chi.Router) { r.Handle("/*", http.StripPrefix("/static", fs)) },
	)
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
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	fmt.Printf("🚀 server is running on port %s ✅\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	log.Fatal(err)
}
