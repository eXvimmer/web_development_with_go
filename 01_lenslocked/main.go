package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles(filepath)
	if err != nil {
		log.Println("error while parsing template file:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Println("error while executing template file:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	executeTemplate(w, filepath.Join("templates", "home.tmpl.html"))
}

func contactHandler(w http.ResponseWriter, _ *http.Request) {
	executeTemplate(w, filepath.Join("templates", "contact.tmpl.html"))
}

func paramHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, chi.URLParam(r, "id"))
}

func faqHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
    <h1>FAQ</h1>
    <h2>Q: Is there a free version?</h2>
    <p>A: Yes, we offer a free trial for 30 days on any paid plans.</p>
    <h2>Q: What are your support hours?</h2>
    <p>A: We have support staff answering emails 24/7, though response times may be a bit slower on weekends.</p>
    <h2>Q: How do I contact support?</h2>
    <p>A: Email us: <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>.</p>
  `)
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
