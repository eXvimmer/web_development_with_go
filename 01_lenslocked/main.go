package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/exvimmer/lenslocked/controllers"
	"github.com/exvimmer/lenslocked/views"
	"github.com/go-chi/chi/v5"
)

type Question struct {
	Text string
	// NOTE: use HTML instead of string, to load things like <a>...</a>. Make
	// sure the source is secure and trusted
	Answer template.HTML
}

var questions = []Question{
	{
		Text:   "Is there a free version?",
		Answer: "Yes, we offer a free trial for 30 days on any paid plans.",
	},
	{
		Text: "What are your support hours?",
		Answer: `We have support staff answering emails 24/7,
    though response times may be a bit slower on weekends.`,
	},
	{
		Text: "How do I contact support?",
		Answer: `Email us: <a href="mailto:support@lenslocked.com">
    support@lenslocked.com</a>.`,
	},
}

func main() {
	r := chi.NewRouter()

	t, err := views.Parse(filepath.Join("templates", "home.tmpl.html"))
	if err != nil {
		panic(err)
	}
	r.Get("/", controllers.StaticHandler(t, nil))

	t, err = views.Parse(filepath.Join("templates", "contact.tmpl.html"))
	if err != nil {
		panic(err)
	}
	r.Get("/contact", controllers.StaticHandler(t, nil))

	t, err = views.Parse(filepath.Join("templates", "faq.tmpl.html"))
	if err != nil {
		panic(err)
	}
	r.Get("/faq", controllers.StaticHandler(t, questions))

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	fmt.Println(" 🚀 server is running on port :3000 ✅")
	err = http.ListenAndServe(":3000", r)
	log.Fatal(err)
}
