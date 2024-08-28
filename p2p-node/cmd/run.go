package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/harryshevan/go-chat-app/p2p-node/p2p"
)

func signalStop() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	fmt.Println("\nGracefully stopped...")

}

func main() {
	target := flag.String("target", "", "target multiaddr string to send a message")
	flag.Parse()

	go p2p.RunNode(*target)
	signalStop()
}
