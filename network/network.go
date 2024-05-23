package network

import (
	"fmt"
	"sync"
)

type Message struct {
	Data string
	From string
}

type Peer struct {
	Address    string
	Connection chan Message
}

var peers = make(map[string]*Peer)
var mu sync.Mutex

func AddPeer(address string, conn chan Message) {
	mu.Lock()
	defer mu.Unlock()
	peers[address] = &Peer{Address: address, Connection: conn}
}

func RemovePeer(address string) {
	mu.Lock()
	defer mu.Unlock()
	delete(peers, address)
}

func GetPeer(address string) (*Peer, bool) {
	mu.Lock()
	defer mu.Unlock()
	peer, exists := peers[address]
	return peer, exists
}

func ListPeers() []string {
	mu.Lock()
	defer mu.Unlock()
	addresses := make([]string, 0, len(peers))
	for address := range peers {
		addresses = append(addresses, address)
	}
	return addresses
}

func SendMessage(address string, msg Message) error {
	peer, exists := GetPeer(address)
	if !exists {
		return fmt.Errorf("peer not found: %s", address)
	}
	peer.Connection <- msg
	return nil
}

func ReceiveMessage(address string) (Message, error) {
	peer, exists := GetPeer(address)
	if !exists {
		return Message{}, fmt.Errorf("peer not found: %s", address)
	}
	msg := <-peer.Connection
	return msg, nil
}

func ConnectToPeer(address string, conn chan Message) {
	AddPeer(address, conn)
}

func HandleConnections() {
	for address, peer := range peers {
		go func(address string, peer *Peer) {
			for msg := range peer.Connection {
				fmt.Printf("Received message from %s: %s\n", address, msg.Data)
			}
		}(address, peer)
	}
}
