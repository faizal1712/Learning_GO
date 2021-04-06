package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

var earnerID, campaignID, companyID int
var serverURL = "http://3.19.86.25:8080"
var socketURL = "3.19.86.25:8081"

// EarnerChat ...
type EarnerChat struct {
	Msg string `json:"newMessage"`
	CID string `json:"CID"`
}

// ClientChat ...
type ClientChat struct {
	Msg string `json:"newMessage"`
	EID string `json:"EID"`
}

// Response ...
type Response struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

// ChatWidgetResponse ...
type ChatWidgetResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    ChatData `json:"data"`
}

// ChatData ...
type ChatData struct {
	ID              int    `json:"id"`
	Host            string `json:"host"`
	HeaderColor     string `json:"headerColor"`
	HeaderFontColor string `json:"headerFontColor"`
	BodyColor       string `json:"bodyColor"`
	BodyFontColor   string `json:"bodyFontColor"`
	FooterColor     string `json:"footerColor"`
	FooterFontColor string `json:"footerFontColor"`
	Token           string `json:"token"`
	CompanyID       int    `json:"companyId"`
}

// MyEventData ...
type MyEventData struct {
	Data  string
	CID   string
	EID   string
	ENAME string
	CDATA ClientData
}

// ClientData ...
type ClientData struct {
	Name       string
	Email      string
	Service    string
	Company    string
	Details    string
	CID        string
	CompanyID  int
	CampaignID int
}

func main() {
	var rooms []string

	var clientList = make(map[string]bool)
	var earnerList = make(map[string]bool)
	var chatting = make(map[string]string)

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel, args interface{}) {
		log.Println("Someone connected")
	})

	err := server.On("CLIENT_WANTS_TO_JOIN", func(c *gosocketio.Channel, cdata ClientData) {
		fmt.Println("Client joined server", cdata.CID)
		room := cdata.Service

		companyID = cdata.CompanyID
		campaignID = cdata.CampaignID

		log.Println("Service selected by client: ", room)
		c.BroadcastTo(room, "CLIENT_JOINED", "Waiting for earner to connect!")
		c.BroadcastTo(room, "CLIENT_WAITING", "Waiting for earner to connect!")

		err := server.On("EARNER_ACCEPT", func(c *gosocketio.Channel, EID string) {
			fmt.Println("Earner joined server", EID)

			channel, _ := server.GetChannel(c.Id())
			err := channel.Emit("EARNER_ACCEPT", cdata)
			if err != nil {
				log.Println("Error>>", err.Error())
			}
			channel, _ = server.GetChannel(cdata.CID)
			err = channel.Emit("EARNER_ACCEPT", EID)
			if err != nil {
				log.Println("Error>>", err.Error())
			}
			chatting[cdata.CID] = EID
			earnerList[EID] = true
			clientList[cdata.CID] = true

			c.BroadcastTo(room, "CLIENT_JOINED", "REQ_ACCEPTED")
			c.BroadcastTo(room, "ack", "Earner Joined! Now you can chat with earner")

			channel, _ = server.GetChannel(EID)
			for _, v := range rooms {
				channel.Leave(string(v))
				fmt.Println("Earner Accpt>>", v, channel.Amount(string(v)))
			}

			fmt.Println(earnerList, "Earner list")
			fmt.Println(clientList, "Client list")
			fmt.Println(chatting, "Chatting")
		})
		if err != nil {
			log.Println("Error>>", err.Error())
		}
	})
	if err != nil {
		log.Println("Error>>", err.Error())
	}

	err = server.On("DATA_FROM_CLIENT", func(c *gosocketio.Channel, data ClientChat) {
		// for k, v := range chatting {
		fmt.Println("Client: ", data.Msg)
		// if k == c.Id() {
		channel, _ := server.GetChannel(data.EID)
		err := channel.Emit("DATA_FROM_CLIENT", MyEventData{Data: data.Msg})
		if err != nil {
			log.Println("Error>>", err.Error())
		}
		// }
		// }
	})
	if err != nil {
		log.Println("Error>>", err.Error())
	}

	err = server.On("DATA_FROM_EARNER", func(c *gosocketio.Channel, data EarnerChat) {
		fmt.Println("dataaaaa")
		// for k, v := range chatting {
		fmt.Println("Earner: ", data.Msg)
		// if v == c.Id() {
		channel, _ := server.GetChannel(data.CID)
		err := channel.Emit("DATA_FROM_EARNER", MyEventData{Data: data.Msg})
		if err != nil {
			log.Println("Error>>", err.Error())
		}
		// }
		// }
	})
	if err != nil {
		log.Println("Error>>", err.Error())
	}

	err = server.On("CLIENT_JOIN", func(c *gosocketio.Channel) {
		clientList[c.Id()] = false
		err := c.Emit("CLIENT_JOIN", MyEventData{Data: "Waiting for earner to connect", CID: c.Id()})
		if err != nil {
			log.Println("Error>>", err.Error())
		}
		fmt.Println(earnerList, "Earner list")
		fmt.Println(clientList, "Client list")
		fmt.Println(chatting, "Chatting")

	})
	if err != nil {
		log.Println("Error>>", err.Error())
	}

	err = server.On("EARNER_JOIN", func(c *gosocketio.Channel, id string) {
		erid, _ := strconv.Atoi(id)
		earnerID = erid
		client := &http.Client{}
		log.Println("Earner room")

		// Create request
		req, err := http.NewRequest("GET", serverURL+"/api/getcompanyservice/"+id, nil)
		if err != nil {
			fmt.Println("errro1: ", err)
			return
		}
		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error2: ", err)
			return
		}
		defer resp.Body.Close()
		// Read Response Body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error>>", err.Error())
		}

		var r Response
		err = json.Unmarshal([]byte(respBody), &r)
		if err != nil {
			log.Println("Error>>", err.Error())
		}
		rooms = r.Data
		log.Println("Earner room", rooms)
		if err != nil {
			fmt.Println("Error3:", err)
			return
		}

		if rooms == nil {
			fmt.Println("Status Offline of Earner")
			return
		}
		for _, v := range rooms {
			c.Join(string(v))
			fmt.Println("Amount>>", v, c.Amount(string(v)))
		}

		earnerList[c.Id()] = false
		err = c.Emit("EARNER_JOIN", MyEventData{Data: "Waiting for client to connect", EID: c.Id()})
		if err != nil {
			log.Println("Error>>", err.Error())
		}
		fmt.Println(earnerList, "Earner list")
		fmt.Println(clientList, "Client list")
		fmt.Println(chatting, "Chatting")
	})
	if err != nil {
		log.Println("Error>>", err.Error())
	}

	err = server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {

		for k, v := range chatting {
			if k == c.Id() {
				channel, _ := server.GetChannel(v)
				err := channel.Emit("CLIENT_LEFT", MyEventData{Data: "Client left the conversation"})
				if err != nil {
					log.Println("Error>>", err.Error())
				}
				// c.Leave("chat")
				delete(chatting, k)
				delete(clientList, k)
				earnerList[v] = false

				payload := map[string]int{"companyId": companyID, "campaignId": campaignID, "earnerId": earnerID}
				jsonValue, _ := json.Marshal(payload)

				// Create request
				client := &http.Client{}
				req, err := http.NewRequest("POST", serverURL+"/api/UpdateEarnerWallet", bytes.NewBuffer(jsonValue))
				if err != nil {
					fmt.Println("errro1: ", err)
					return
				}
				// Fetch Request
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println("Error2: ", err)
					return
				}
				defer resp.Body.Close()
				// Read Response Body
				respBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println("Error>>", err.Error())
				}
				log.Println("response wallet body", string(respBody))
				break
			}
			if v == c.Id() {
				channel, _ := server.GetChannel(k)
				err := channel.Emit("EARNER_LEFT", MyEventData{Data: "Earner left the conversation"})
				if err != nil {
					log.Println("Error>>", err.Error())
				}
				// c.Leave("chat")
				delete(chatting, k)
				delete(earnerList, v)
				clientList[k] = false
				break
			}
		}

		for k := range earnerList {
			if k == c.Id() {
				// c.Leave("chat")
				delete(earnerList, k)
				break
			}
		}

		for k := range clientList {
			if k == c.Id() {
				// c.Leave("chat")
				delete(clientList, k)
				break
			}
		}

		fmt.Println("Chatting", chatting)
		fmt.Println("Clients", clientList)
		fmt.Println("Earners", earnerList)
	})
	if err != nil {
		log.Println("Error>>", err.Error())
	}

	http.Handle("/", Middleware(http.FileServer(http.Dir("../client_chat")), func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

			u, err := url.Parse(req.Header.Get("Referer"))
			if err != nil {
				log.Fatal(err)
			}
			// log.Println("HOST", u.Host)
			client := &http.Client{}

			// Create Request
			request, err := http.NewRequest("GET", serverURL+"/api/chatsubscription/"+u.Host, nil)
			if err != nil {
				fmt.Println("Error In Create Request: ", err)
				return
			}
			// Fetch Request
			resp, err := client.Do(request)
			if err != nil {
				fmt.Println("Error In Fetch Request: ", err)
				return
			}
			defer resp.Body.Close()
			// Read Response Body
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error>>", err.Error())
			}

			var r ChatWidgetResponse
			err = json.Unmarshal([]byte(respBody), &r)
			if err != nil {
				log.Println("Error>>", err.Error())
			}
			// log.Println("Heelllllllo", r.Data.Host)
			if u.Host == r.Data.Host || u.Host == socketURL {
				// log.Println("inside3", req.Header.Get("Referer"))
				// log.Println("1", u.Host)
				next.ServeHTTP(res, req)
			}
		})
	},
	))
	http.Handle("/socket.io/", server)
	// http.Handle("/", http.FileServer(http.Dir("../client_chat")))

	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}

}

// Middleware ...
func Middleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}
