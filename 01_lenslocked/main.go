package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

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

func executeTemplate(w http.ResponseWriter, filepath string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles(filepath)
	if err != nil {
		log.Println("error while parsing template file:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Println("error while executing template file:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	executeTemplate(w, filepath.Join("templates", "home.tmpl.html"), nil)
}

func contactHandler(w http.ResponseWriter, _ *http.Request) {
	executeTemplate(w, filepath.Join("templates", "contact.tmpl.html"), nil)
}

func paramHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, chi.URLParam(r, "id"))
}

func faqHandler(w http.ResponseWriter, _ *http.Request) {
	executeTemplate(w, filepath.Join("templates", "faq.tmpl.html"), questions)
}

func main() {
	r := chi.NewRouter()

	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.Get("/{id}", paramHandler)
	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	fmt.Println(" 🚀 server is running on port :3000 ✅")
	err := http.ListenAndServe(":3000", r)
	log.Fatal(err)
}
