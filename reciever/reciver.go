package main

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true;
	},
}

var connections = []*websocket.Conn{}

func remove(
	s []*websocket.Conn,
	i int,
) []*websocket.Conn {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	connections = append(connections, conn)

	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	for {
		_, mess, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			for i, c := range connections {
				if c == conn {
					connections = remove(connections, i)
					log.Println("deleted connection")
					break
				}
			}
			break
		}

		log.Printf("Received message: %s", mess)

		for _, c := range connections {
			if c == conn {
				continue
			}

			err = c.WriteMessage(websocket.TextMessage, mess)

			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func startWebSocket() {
	http.HandleFunc("/websocket", websocketHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	startWebSocket()
}
