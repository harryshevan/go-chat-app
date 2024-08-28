package p2p

import (
	"fmt"
	"log"

	host "github.com/libp2p/go-libp2p/core/host"
	peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type DiscoveryNotifee struct {
	h        host.Host
	peerChan chan peer.AddrInfo
}

// HandlePeerFound implements the HandlePeerFound method of the Notifee interface
func (n *DiscoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("Discovered a new peer: %s %s\n", pi.ID.Loggable(), pi.Addrs)
	n.peerChan <- pi
	// if err := n.h.Connect(context.Background(), pi); err != nil {
	// 	fmt.Printf("Error connecting to new peer %s: %s\n", pi.ID.Loggable(), err)
	// } else {
	// 	fmt.Printf("Connected to new peer: %s\n", pi.ID.Loggable())
	// }
}

func SetupMDNS(h host.Host, serviceTag string, peerChan chan peer.AddrInfo) {
	notifee := &DiscoveryNotifee{h: h, peerChan: peerChan}
	service := mdns.NewMdnsService(h, serviceTag, notifee)
	if err := service.Start(); err != nil {
		log.Fatalf("mDNS service failed to start: %s", err)
	}
	fmt.Println("mDNS service started")
}
