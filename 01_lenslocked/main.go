package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/exvimmer/lenslocked/controllers"
	"github.com/exvimmer/lenslocked/templates"
	"github.com/exvimmer/lenslocked/views"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgressConfig struct {
	host, port, user, password, dbname, sslmode string
}

func (p *PostgressConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.host, p.port, p.user, p.password, p.dbname, p.sslmode)
}

func main() {
	openDB()
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	r.Route("/static", func(r chi.Router) {
		r.Handle("/*", http.StripPrefix("/static", fs))
	})

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

	usersC := controllers.Users{
		Templates: controllers.UsersTemplates{
			New: views.Must(views.ParseFS(templates.FS, "signup.tmpl.html",
				"tailwind.tmpl.html")),
		},
	}
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	fmt.Println(" 🚀 server is running on port :3000 ✅")
	err := http.ListenAndServe(":3000", r)
	log.Fatal(err)
}

func openDB() {
	pgConfig := PostgressConfig{
		host:     "localhost",
		port:     "5432",
		user:     "goblina",
		password: "jinnythejimbo",
		dbname:   "lenslocked",
		sslmode:  "disable",
	}
	db, err := sql.Open("pgx", pgConfig.String())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("connected")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			first_name TEXT,
			last_name TEXT,
			email TEXT UNIQUE NOT NULL,
			age INT
		);

		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id),
			amount INT,
			description TEXT
		);
`)
	if err == nil {
		fmt.Println("tables are ready")
	} else {
		panic(err)
	}
}
