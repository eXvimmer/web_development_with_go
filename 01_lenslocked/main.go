package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my great web site!</h1>")
}

func contactHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>Contact Page</h1><p>To get in touch, send me an
  <a href="mailto:mustafa.hayati1992@gmail.com">email</a>.</p>`)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/contact", contactHandler)
	fmt.Println(" 🚀 server is running on port :3000 ✅")
	http.ListenAndServe(":3000", nil)
}
