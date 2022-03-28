package main

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var purs = []string{
	"pur",
	"мур",
	"veh",
	"вех",
	"мурище",
	"vehbot",
	"мурмурмур",
	"муррррррр",
	"vehvehveh",
	"vehhhh",
	"мурк",
	"мурена",
	"самый мурный мур из всех мурных муров",
}

type Client struct {
	Name     uuid.UUID
	Messages []string
}

var clients = make(map[*websocket.Conn]*Client)

var pur = false

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	go echoPur()
	http.ListenAndServe(":8080", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(purs)
	val := purs[n]
	err = ws.WriteMessage(1, []byte(val))
	if err != nil {
		log.Println(err)
	}
	id := uuid.New()
	clients[ws] = &Client{id, make([]string, 0)}
	reader(ws)

}

func echoPur() {
	rand.Seed(time.Now().Unix())
	for {
		if pur {
			n := rand.Int() % len(purs)
			val := purs[n]
			if rand.Intn(2) == 1 {
				val = strings.Title(val)
			}
			time.Sleep(time.Millisecond * 800)
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, []byte(val))
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
		command := strings.Split(string(p), " ")
		switch command[0] {
		case "ls":
			for _, v := range clients {
				conn.WriteMessage(websocket.TextMessage, []byte(v.Name.String()))
			}

		case "i":
			conn.WriteMessage(websocket.TextMessage, []byte(clients[conn].Name.String()))

		case "read":
			readCase(command, conn)

		case "write":
			writeCase(command, conn)

		case "рыжая":
			pur = true

		default:
			conn.WriteMessage(websocket.TextMessage, []byte("Unknown command"))
		}

	}
}

func writeCase(command []string, conn *websocket.Conn) {
	if len(command) > 2 {
		for _, v := range clients {
			if v.Name.String() == command[1] {
				v.Messages = append(v.Messages, strings.Join(command[2:], " "))
				log.Println(v.Messages)
				return
			}
		}

		conn.WriteMessage(websocket.TextMessage, []byte("Unknown name"))

	} else {
		conn.WriteMessage(websocket.TextMessage, []byte("Bad request"))
	}
}

func readCase(command []string, conn *websocket.Conn) {
	if len(command) > 1 {

		for _, v := range clients {
			if v.Name.String() == command[1] {
				log.Println(v.Messages)
				conn.WriteMessage(websocket.TextMessage, []byte(strings.Join(v.Messages, "; ")))
				return
			}
		}

		conn.WriteMessage(websocket.TextMessage, []byte("Unknown name"))

	} else {
		conn.WriteMessage(websocket.TextMessage, []byte(strings.Join(clients[conn].Messages, "; ")))
	}
}
