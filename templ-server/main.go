package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

type Person struct {
	Name string
}

func main() {
	p := Person{Name: "kek"}
	component := Greeting(p)

	http.Handle("/", templ.Handler(component))
	http.Handle("/404", templ.Handler(notFoundComponent(), templ.WithStatus(http.StatusNotFound)))

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
