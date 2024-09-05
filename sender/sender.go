package main

import (
	"log"
	"github.com/gorilla/websocket"
	"bufio"
	"os"
	"net/url"
	"fmt"
)

func readMessages(conn *websocket.Conn) {
	for {
		_, mess, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			break
		}

		fmt.Printf("sender: %s\n", mess)
	}
}

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/websocket"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	log.Printf("write a message")
  
	go readMessages(conn)

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		err := scanner.Err()
    if err != nil {
        log.Fatal(err)
				return
    }

		text := scanner.Bytes()

		err = conn.WriteMessage(websocket.TextMessage, text)

		if err != nil {
			log.Println(err)
			return
		}
	}
}
