package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my awesome web site!</h1>")
}

func main() {
	http.HandleFunc("/", homeHandler)
	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", nil)
}
