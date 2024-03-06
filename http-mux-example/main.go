package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "Wildcard test (path: %q). ID is %s", html.EscapeString(r.URL.Path), id)
	}

	http.HandleFunc("/hello", handler)

	http.HandleFunc("/hello/{id}", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
