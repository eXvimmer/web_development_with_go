package main

import (
	"html/template"
	"os"
)

type UserMeta struct {
	Visited int
}

type User struct {
	Name string
	Age  int
	Meta UserMeta
}

func main() {
	t, err := template.ParseFiles("hello.tmpl.html")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "Mustafa Hayati",
		Age:  30,
		Meta: UserMeta{Visited: 4},
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
