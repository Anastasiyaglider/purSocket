package main

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

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
	"самый мурный мур из всех мурных  муров",
}

var clients = make(map[*websocket.Conn]bool)

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
	clients[ws] = true
}

func echoPur() {
	rand.Seed(time.Now().Unix())
	for {
		n := rand.Int() % len(purs)
		val := purs[n]
		if rand.Intn(2) == 1 {
			val = strings.Title(val)
		}
		time.Sleep(time.Millisecond * 10)
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
