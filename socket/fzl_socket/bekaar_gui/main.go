package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type Channel struct {
	Channel string `json:"channel"`
}

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

var router *mux.Router

func homeHandler(tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r)
	})
}

func main() {
	tpl := template.Must(template.ParseFiles("jaimeen_index.html"))
	router = mux.NewRouter()
	router.Handle("/home", homeHandler(tpl))

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	router.Handle("/socket.io/", server)
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("Created lobby")
		log.Println("New client connected, client id is ", c.Id())
		c.Join("Room111111111111111111111")
		channel, _ := server.GetChannel(c.Id())
		channel.Emit("/message", Message{10, "main", "using emit"})
		fmt.Println("------------------------------------")
		server.BroadcastTo("Room111111111111111111111", "/message", Message{10, "main", "using broadcast"})
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("Disconnected")
	})

	server.On("/join", func(c *gosocketio.Channel, channel Channel) string {
		log.Println("Client joined to ", channel.Channel)
		return "joined to " + channel.Channel
	})

	server.On("/connect", func(c *gosocketio.Channel, channel Channel) string {
		log.Println("Client joined to ", channel.Channel)
		return "joined to " + channel.Channel
	})

	log.Fatal(http.ListenAndServe(":8080", router))

	fmt.Println("")
}
