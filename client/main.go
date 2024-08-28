package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	viper.BindEnv("WS_SERVER_URL")
	url := viper.Get("WS_SERVER_URL")
	fmt.Println("Connecting to the server at:", url)

	// os.Exit(0)
	wsURL := url.(string)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("Dial:", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	// Reading messages from the server
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read:", err)
				return
			}
			log.Printf("Received: %s", message)
		}
	}()

	// Sending a message to the server
	err = conn.WriteMessage(websocket.TextMessage, []byte("kek puk!"))
	if err != nil {
		log.Println("Write:", err)
		return
	}

	// Signal handling to close the connection
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupt received, closing connection")

			// Close the WebSocket connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
