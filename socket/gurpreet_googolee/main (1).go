package main

import (
	"fmt"
	"log"
	"net/http"

	guuid "github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
)

func main() {
	server, serveError := socketio.NewServer(nil)
	if serveError != nil {
		log.Fatalln(serveError)
	}

	// server.OnConnect("/", func(s socketio.Conn) error {
	// 	s.SetContext("")
	// 	fmt.Println("connected:", s.ID())
	// 	return nil

	// })

	//Add all connected user to a room, in example? "bcast"
	server.OnConnect("/", func(s socketio.Conn) error {

		s.SetContext("")
		fmt.Println("connected:", s.ID())
		//s.Join("bcast")
		return nil
	})

	//Broadcast message to all connected user

	// server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
	// 	fmt.Println("Messaage from Client:", msg)

	// s.Emit("reply", "MEsssage from server "+msg)

	// 	//Sending one client data to all that are connected in room

	// })

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg, roomname string) string {
		s.SetContext(msg)
		fmt.Println(msg)
		fmt.Println(roomname)
		//s.Join("bcast")
		//Acknowledgement from server to all clients that are in room
		//server.BroadcastToRoom("", roomname, "reply", msg)
		server.BroadcastToRoom("", "bcast", "reply", msg)
		//Acknowledgement from server to client
		s.Emit("reply", "MEsssage from server "+msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	server.OnEvent("/roomcreate", "roomcreate", func(s socketio.Conn) string {
		roomname := guuid.New().String()
		s.Join(roomname)
		fmt.Println("roomname ", roomname)
		// s.Emit("roomname", roomname)
		return roomname
	})

	server.OnEvent("/roomjoin", "roomjoin", func(s socketio.Conn, roomname string) {
		s.Join(roomname)
		fmt.Println("join the room ", roomname)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
