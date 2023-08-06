package main

import (
	"html/template"
	"os"
)

type UserMeta struct {
	Visited int
}

type User struct {
	Name         string
	SafeBio      string
	DangerousBio template.HTML // NOTE: this should come from a trusted source
}

func main() {
	t, err := template.ParseFiles("hello.tmpl.html")
	if err != nil {
		panic(err)
	}

	user := User{
		Name:         "Mustafa Hayati",
		SafeBio:      `<script>alert("hi")</script>`,
		DangerousBio: `<script>alert("hi")</script>`,
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
