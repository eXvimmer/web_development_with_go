package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/exvimmer/lenslocked/controllers"
	"github.com/exvimmer/lenslocked/templates"
	"github.com/exvimmer/lenslocked/views"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	r.Route("/static", func(r chi.Router) {
		r.Handle("/*", http.StripPrefix("/static", fs))
	})

	r.Get("/",
		controllers.StaticHandler(
			views.Must(views.ParseFS(templates.FS, "home.tmpl.html", "tailwind.tmpl.html")),
			nil))

	r.Get("/contact",
		controllers.StaticHandler(
			views.Must(views.ParseFS(templates.FS, "contact.tmpl.html", "tailwind.tmpl.html")),
			nil))

	r.Get("/faq",
		controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.tmpl.html", "tailwind.tmpl.html"))))

	r.Get("/signup",
		controllers.StaticHandler(
			views.Must(views.ParseFS(templates.FS, "signup.tmpl.html", "tailwind.tmpl.html")),
			nil))

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	fmt.Println(" 🚀 server is running on port :3000 ✅")
	err := http.ListenAndServe(":3000", r)
	log.Fatal(err)
}
