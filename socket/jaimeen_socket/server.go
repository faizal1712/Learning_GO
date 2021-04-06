package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var router *mux.Router
var Listhub []*hub

func homeHandler(tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r)
	})
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("home.html"))
	tpl.Execute(w, nil)
}

func main() {
	tpl := template.Must(template.ParseFiles("index.html"))
	router = mux.NewRouter()
	router.Handle("/", homeHandler(tpl))
	router.HandleFunc("/ws/create", createLobby)
	router.HandleFunc("/ws/join/{id}", joinLobby)
	router.HandleFunc("/home", HomePage)
	log.Printf("serving on port http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

//create len() = 12, 11 listhub[11]
