package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my great web site!</h1>")
}

func main() {
	http.HandleFunc("/", homeHandler)
	fmt.Println(" 🚀 server is running on port :3000 ✅")
	http.ListenAndServe(":3000", nil)
}
