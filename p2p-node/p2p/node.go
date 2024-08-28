package p2p

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p/core/host"
	network "github.com/libp2p/go-libp2p/core/network"
	peer "github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

const protocolID = "/simple/1.0.0"
const responseMessage = "Hello from the node!"

func RunNode(target string) {
	h, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	// enable local discovery
	peerChan := make(chan peer.AddrInfo)
	go SetupMDNS(h, "local-discovery", peerChan)

	if target == "" {
		// Start listening for incoming connections
		h.SetStreamHandler(protocolID, handleStream)
		fmt.Printf("Listening on addresses:\n")
		for _, addr := range h.Addrs() {
			fmt.Printf(" - %s/p2p/%s\n", addr, h.ID().Loggable())
		}

		// Prevent the application from exiting
		select {}
	} else {
		k := <-peerChan
		sendMessageByAddr(h, k)
	}

}

func handleStream(s network.Stream) {
	fmt.Println("Received a new stream!")

	reader := bufio.NewReader(s)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Stream closed:", err)
			return
		}
		fmt.Print("Received:", msg)

		_, err = s.Write([]byte(responseMessage + "\n"))
		if err != nil {
			fmt.Println("Error writing to stream:", err)
			return
		}
	}
}

func userInp(s network.Stream) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter some text (type 'exit' to quit): ")
		if scanner.Scan() {
			input := scanner.Text()
			if strings.ToLower(input) == "exit" {
				fmt.Println("Exiting...")
				break
			}
			reader := bufio.NewReader(s)
			reply, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading from stream:", err)
				return
			}

			fmt.Println("Received reply:", strings.TrimSpace(reply))
		} else {
			// Handle scanner errors if any
			fmt.Println("Error reading input:", scanner.Err())
			break
		}
	}

}

func sendMessageByStr(h host.Host, target string) {
	maddr, err := ma.NewMultiaddr(target)
	if err != nil {
		panic(err)
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := h.Connect(ctx, *info); err != nil {
		fmt.Println("Connection failed:", err)
		return
	}
	fmt.Println("Connected to", target)

	s, err := h.NewStream(ctx, info.ID, protocolID)
	if err != nil {
		panic(err)
	}

	userInp(s)
}

func sendMessageByAddr(h host.Host, target peer.AddrInfo) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := h.Connect(ctx, target); err != nil {
		fmt.Println("Connection failed:", err)
		return
	}
	fmt.Println("Connected to", target)

	s, err := h.NewStream(ctx, target.ID, protocolID)
	if err != nil {
		panic(err)
	}

	userInp(s)
}
