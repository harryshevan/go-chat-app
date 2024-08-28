package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p/core/host"
	peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// DiscoveryNotifee defines a structure that implements the Notifee interface
type DiscoveryNotifee struct {
	h        host.Host
	peerChan chan peer.AddrInfo
}

// HandlePeerFound implements the HandlePeerFound method of the Notifee interface
func (n *DiscoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("Discovered a new peer: %s\n", pi.ID.Loggable())
	n.peerChan <- pi
	if err := n.h.Connect(context.Background(), pi); err != nil {
		fmt.Printf("Error connecting to new peer %s: %s\n", pi.ID.Loggable(), err)
	} else {
		fmt.Printf("Connected to new peer: %s\n", pi.ID.Loggable())
	}
}

func setupMDNS(h host.Host, serviceTag string, peerChan chan peer.AddrInfo) {
	notifee := &DiscoveryNotifee{h: h, peerChan: peerChan}
	service := mdns.NewMdnsService(h, serviceTag, notifee)
	if err := service.Start(); err != nil {
		log.Fatalf("mDNS service failed to start: %s", err)
	}
	fmt.Println("mDNS service started")
}

func main() {
	// Create a new libp2p host with default options
	h, err := libp2p.New()
	if err != nil {
		log.Fatalf("Failed to create libp2p host: %s", err)
	}

	fmt.Printf("Node ID: %s\n", h.ID().Loggable())
	for _, addr := range h.Addrs() {
		fmt.Printf("Address: %s\n", addr.String())
	}

	// Channel to collect discovered peers
	peerChan := make(chan peer.AddrInfo)

	// Set up mDNS discovery
	go setupMDNS(h, "local-discovery", peerChan)

	// Keep track of discovered peers
	discoveredPeers := make(map[peer.ID]peer.AddrInfo)

	// Run the discovery for 30 seconds
	timer := time.NewTimer(30 * time.Second)
	for {
		select {
		case pi := <-peerChan:
			discoveredPeers[pi.ID] = pi
		case <-timer.C:
			fmt.Println("Discovery finished. List of discovered peers:")
			for id, info := range discoveredPeers {
				fmt.Printf("Peer ID: %s, Addresses: %s\n", id.Loggable(), info.Addrs)
			}
			return
		}
	}
}
