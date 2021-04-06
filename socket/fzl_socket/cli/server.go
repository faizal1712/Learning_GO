package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/17twenty/shortuuid"
	"github.com/gorilla/mux"

	guuid "github.com/google/uuid"
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

func main() {
	router = mux.NewRouter()
	router.Handle("/", http.FileServer(http.Dir("./asset")))
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	// server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
	// log.Println("Created lobby", counter)

	// c.Emit("/message", Message{10, "main", "using emit"})
	// var id string
	// id = str[counter]
	// c.Join(id)
	// c.BroadcastTo(id, "/message", Message{10, "main", "using broadcast"})
	// counter++
	// })

	server.On("/create", func(c *gosocketio.Channel) string {

		c.Emit("/message", Message{10, "main", "using emit"})
		_ = guuid.New().String()
		id := shortuuid.New()
		if len(c.List(id)) > 0 {
			fmt.Println("Error creating code")
			return "Error Creating Code"
		}
		fmt.Println(id)
		c.Join(id)
		log.Println("Created lobby", id)
		fmt.Println(c.List(id))
		c.BroadcastTo(id, "/message", Message{10, "main", "using broadcast"})
		return id
	})

	server.On("/message", func(c *gosocketio.Channel, args Message) {
		fmt.Println("receiving message from server")
		chanee, _ := server.GetChannel(c.Id())
		fmt.Println(chanee.Id(), c.Id())
		c.BroadcastTo(c.Id(), "/message", Message{10, "main", "using broadcast"})
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("Disconnected")
	})

	server.On("/join", func(c *gosocketio.Channel, channel Channel) string {
		if len(c.List(channel.Channel)) == 0 {
			return "Need to create a channel " + channel.Channel
		}
		log.Println("Client joined to ", channel.Channel)
		c.Join(channel.Channel)
		c.BroadcastTo(channel.Channel, "/message", Message{10, "main", "using broadcast"})
		return "joined to " + channel.Channel
	})
	router.Handle("/socket.io/", server)

	log.Fatal(http.ListenAndServe(":8080", router))

	fmt.Println("")
}

// "/message/id"
