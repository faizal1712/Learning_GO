package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
		// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.ListenAndServe(":80", nil)

	// h1 := func(w http.ResponseWriter, _ *http.Request) {
	// 	io.WriteString(w, "Hello from a HandleFunc #1!\n")
	// }
	// h2 := func(w http.ResponseWriter, _ *http.Request) {
	// 	io.WriteString(w, "Hello from a HandleFunc #2!\n")
	// }

	// http.HandleFunc("/", h1)
	// http.HandleFunc("/endpoint", h2)
}
