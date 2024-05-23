package main

import (
	"blockchain-hello-golang/block"
	"blockchain-hello-golang/network"
	"blockchain-hello-golang/transaction"
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
)

var nodes = make(map[int]*Node)
var mu sync.Mutex
var minedBlocks = 0

const blockTarget = 4

type Node struct {
	id             int
	color          string
	mining         bool
	peers          map[int]*Node
	txChannel      chan transaction.Transaction
	blockChannel   chan block.Block
	messageChannel chan network.Message
	chain          []block.Block
}

func newNode(id int, color string, mining bool) *Node {
	return &Node{
		id:             id,
		color:          color,
		mining:         mining,
		peers:          make(map[int]*Node),
		txChannel:      make(chan transaction.Transaction, 100),
		blockChannel:   make(chan block.Block, 100),
		messageChannel: make(chan network.Message, 100),
	}
}

func (n *Node) log(action string) {
	log.Printf("%s[node-%d] %s%s\n", n.color, n.id, action, resetColor)
}

func (n *Node) discoverPeers() {
	for {
		mu.Lock()
		if len(n.peers) < avgPeers {
			peerID := rand.Intn(len(nodes))
			if peerID != n.id && len(nodes[peerID].peers) < maxPeers {
				n.peers[peerID] = nodes[peerID]
				nodes[peerID].peers[n.id] = n
				n.log(fmt.Sprintf("Connected to peer %d", peerID))
			}
		}
		mu.Unlock()
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
}

func (n *Node) dropPeers() {
	for {
		mu.Lock()
		if len(n.peers) > avgPeers {
			for peerID := range n.peers {
				delete(n.peers, peerID)
				delete(nodes[peerID].peers, n.id)
				n.log(fmt.Sprintf("Dropped peer %d", peerID))
				break
			}
		}
		mu.Unlock()
		time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
	}
}

func (n *Node) sendMessage(to int, msg network.Message) {
	if peer, ok := n.peers[to]; ok {
		msg.From = strconv.Itoa(n.id) // Set the From field
		peer.messageChannel <- msg
	}
}

func (n *Node) receiveMessages() {
	for msg := range n.messageChannel {
		n.log(fmt.Sprintf("Received message from %s: %s", msg.From, msg.Data))
	}
}

func (n *Node) generateTransactions() {
	for {
		to := rand.Intn(len(nodes))
		if to != n.id {
			txID := fmt.Sprintf("tx-%d-%d", n.id, rand.Int())
			tx := transaction.Transaction{ID: txID, Inputs: []transaction.Input{}, Outputs: []transaction.Output{{Value: rand.Intn(100), ScriptPubKey: fmt.Sprintf("address-%d", to)}}}
			n.txChannel <- tx
			n.log(fmt.Sprintf("Generated transaction to %d: %+v", to, tx))
		}
		time.Sleep(time.Duration(rand.Intn(20)+10) * time.Second)
	}
}

func (n *Node) handleTransactions() {
	for tx := range n.txChannel {
		n.log(fmt.Sprintf("Handling transaction: %+v", tx))
	}
}

func (n *Node) mineBlock() {
	for {
		if !n.mining {
			return
		}
		time.Sleep(time.Duration(rand.Intn(30)+10) * time.Second)
		mu.Lock()
		if minedBlocks >= blockTarget {
			mu.Unlock()
			return
		}
		previousHash := "0"
		if len(n.chain) > 0 {
			previousHash = n.chain[len(n.chain)-1].Hash
		}
		timestamp := time.Now().Unix()
		transactions := []transaction.Transaction{}
		for i := 0; i < rand.Intn(5)+1; i++ {
			to := rand.Intn(len(nodes))
			if to != n.id {
				txID := fmt.Sprintf("tx-%d-%d", n.id, rand.Int())
				transactions = append(transactions, transaction.Transaction{ID: txID, Inputs: []transaction.Input{}, Outputs: []transaction.Output{{Value: rand.Intn(100), ScriptPubKey: fmt.Sprintf("address-%d", to)}}})
			}
		}
		nonce := 0
		var hash string
		for {
			data := strconv.Itoa(len(n.chain)) + previousHash + strconv.FormatInt(timestamp, 10) + fmt.Sprintf("%v", transactions) + strconv.Itoa(nonce)
			hashBytes := sha256.Sum256([]byte(data))
			hash = fmt.Sprintf("%x", hashBytes)
			if hash[:miningDifficulty] == "0000" {
				break
			}
			nonce++
		}
		blk := block.Block{Index: len(n.chain), PrevHash: previousHash, Timestamp: timestamp, Transactions: transactions, Hash: hash, Nonce: nonce, Difficulty: miningDifficulty}
		n.chain = append(n.chain, blk)
		n.log(fmt.Sprintf("Mined block: %+v", blk))
		minedBlocks++
		for _, peer := range n.peers {
			peer.blockChannel <- blk
		}
		mu.Unlock()
		if minedBlocks >= blockTarget {
			return
		}
	}
}

func (n *Node) receiveBlocks() {
	for blk := range n.blockChannel {
		mu.Lock()
		valid := true
		if len(n.chain) > 0 && n.chain[len(n.chain)-1].Hash != blk.PrevHash {
			valid = false
		}
		for _, b := range n.chain {
			if b.Hash == blk.Hash {
				valid = false
				break
			}
		}
		if valid {
			n.chain = append(n.chain, blk)
			n.log(fmt.Sprintf("Accepted block: %+v", blk))
		} else {
			n.log(fmt.Sprintf("Rejected block: %+v", blk))
		}
		mu.Unlock()
	}
}

func logNetworkStructure() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		mu.Lock()
		fmt.Println("Network Structure:")
		nodeIDs := make([]int, 0, len(nodes))
		for id := range nodes {
			nodeIDs = append(nodeIDs, id)
		}
		sort.Ints(nodeIDs)
		for _, id := range nodeIDs {
			node := nodes[id]
			peerIDs := []int{}
			for peerID := range node.peers {
				peerIDs = append(peerIDs, peerID)
			}
			sort.Ints(peerIDs)
			fmt.Printf("node %d -> %v\n", id, peerIDs)
			fmt.Printf("Transactions: %v\n", getTransactions(node.txChannel))
		}
		fmt.Println()
		mu.Unlock()
	}
}

func getTransactions(txChannel chan transaction.Transaction) []transaction.Transaction {
	var transactions []transaction.Transaction
	for {
		select {
		case tx := <-txChannel:
			transactions = append(transactions, tx)
		default:
			return transactions
		}
	}
}

var colors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Magenta
	"\033[36m", // Cyan
}

const resetColor = "\033[0m"
const miningDifficulty = 4
const maxPeers = 5
const avgPeers = 3

func main() {
	rand.Seed(time.Now().UnixNano())

	// Initialize nodes
	for i := 0; i < 100; i++ {
		color := colors[i%len(colors)]
		mining := rand.Float32() < 0.5 // Randomly assign mining capability to half the nodes
		nodes[i] = newNode(i, color, mining)
	}

	// Start the network
	for _, node := range nodes {
		go node.discoverPeers()
		go node.dropPeers()
		go node.receiveMessages()
		go node.generateTransactions()
		go node.handleTransactions()
		go node.mineBlock()
		go node.receiveBlocks()
	}

	// Log network structure periodically
	go logNetworkStructure()

	// Run until the target number of blocks are mined
	for {
		mu.Lock()
		if minedBlocks >= blockTarget {
			mu.Unlock()
			break
		}
		mu.Unlock()
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Blockchain simulation complete.")
}