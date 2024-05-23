package peer

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Peer struct {
	Address    string
	Connection net.Conn
}

var peers = make(map[string]*Peer)
var mu sync.Mutex

func AddPeer(address string, conn net.Conn) {
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

func ConnectToPeer(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	AddPeer(address, conn)
	return conn, nil
}

func ListenForPeers(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		AddPeer(conn.RemoteAddr().String(), conn)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	address := conn.RemoteAddr().String()
	defer RemovePeer(address)

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error reading from connection:", err)
			return
		}
		log.Printf("Received message from %s: %s\n", address, string(buffer[:n]))
		processMessage(buffer[:n])
	}
}

func processMessage(msg []byte) {
	// Implement message parsing and handling logic here
	log.Printf("Processing message: %s\n", msg)
}

func sendHeartbeat(conn net.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		_, err := conn.Write([]byte("ping"))
		if err != nil {
			log.Println("Error sending heartbeat:", err)
			return
		}
	}
}
