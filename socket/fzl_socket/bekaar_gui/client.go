// package main

// import (
// 	"flag"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/gorilla/mux"
// 	gosocketio "github.com/graarh/golang-socketio"
// 	"github.com/graarh/golang-socketio/transport"
// )

// type Channel struct {
// 	Channel string `json:"channel"`
// }

// type Message struct {
// 	Id      int    `json:"id"`
// 	Channel string `json:"channel"`
// 	Text    string `json:"text"`
// }

// var router *mux.Router

// var Cl *gosocketio.Client

// func main() {
// 	fmt.Println("")
// 	router = mux.NewRouter()
// 	create := flag.Bool("create", false, "a bool")
// 	join := flag.String("join", "abcd", "a string")
// 	flag.Parse()
// 	Cl, err := gosocketio.Dial(
// 		gosocketio.GetUrl("localhost", 8080, false),
// 		transport.GetDefaultWebsocketTransport())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if *create {
// 		createLobby(Cl)
// 	} else {
// 		joinLobby(*join, Cl)
// 	}

// 	err = Cl.On("/message", func(h *gosocketio.Channel, args Message) {
// 		log.Println("--- Got chat message: ", args)
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	time.Sleep(60000000 * time.Second)
// }

// func createLobby(Cl *gosocketio.Client) {
// 	err := Cl.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
// 		log.Println("Connected")
// 	})
// 	fmt.Println("create")
// 	if err != nil {
// 		log.Fatal(err, "createlobby")
// 	}
// }

// func joinLobby(channel string, Cl *gosocketio.Client) {
// 	log.Println("Acking /join")
// 	result, err := Cl.Ack("/join", Channel{channel}, time.Second*5)
// 	if err != nil {
// 		log.Fatal(err)
// 	} else {
// 		log.Println("Ack result to /join: ", result)
// 	}
// }
