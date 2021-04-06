package main

import (
	"flag"
	"fmt"
	"log"

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

var Cl *gosocketio.Client

func main() {
	fmt.Println("")
	router = mux.NewRouter()
	create := flag.Bool("create", false, "a bool")
	join := flag.String("join", "", "a string")
	flag.Parse()
	Cl, err := gosocketio.Dial(
		gosocketio.GetUrl("localhost", 8080, false),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}
	if *create {
		fmt.Println("if")
		createLobby(Cl)
	} else {
		fmt.Println("else")
		joinLobby(*join, Cl)
	}

	err = Cl.On("/message", func(h *gosocketio.Channel, args Message) {
		log.Println("--- Got chat message: ", args)
	})
	if err != nil {
		log.Fatal(err)
	}
	// time.Sleep(60000000 * time.Second)
	str := ""
	for str != "bye" {
		fmt.Print("Enter your message :")
		fmt.Scanf("%s", &str)
		sendMessage(Cl, str)
	}
}

func sendMessage(Cl *gosocketio.Client, str string) {
	fmt.Println("sending message")
	err := Cl.Emit("/message", Message{1, "", "Hello"})
	if err != nil {
		fmt.Println(err)
	}
}

func createLobby(Cl *gosocketio.Client) {
	err := Cl.Emit("/create", nil)
	fmt.Println("create")
	if err != nil {
		log.Fatal(err, "createlobby")
	}
}

func joinLobby(channel string, Cl *gosocketio.Client) {
	log.Println("Acking /join")
	err := Cl.Emit("/join", Channel{channel})
	// result, err := Cl.Ack("/join", Channel{channel}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Ack result to /join: ")
		// log.Println("Ack result to /join: ", result)
	}
}
