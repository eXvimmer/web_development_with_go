package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my awesome web site!</h1>")
}

func contactHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>Contact Page</h1><p>To get in touch, send me an
  <a href="mailto:mustafa.hayati1992@gmail.com">email</a>.</p>`)
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

type Router struct{}

func (Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	default:
		http.Error(w,
			http.StatusText(http.StatusNotFound),
			http.StatusNotFound,
		)
	}
}

func main() {
	fmt.Println(" 🚀 server is running on port :3000 ✅")
	http.ListenAndServe(":3000", Router{})
}
