package utxo

import (
	"blockchain-hello-golang/network"
	"blockchain-hello-golang/transaction"
	"sync"
)

var utxoSet = make(map[string]map[int]transaction.Output)
var mu sync.Mutex

func addUTXO(tx transaction.Transaction) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := utxoSet[tx.ID]; !exists {
		utxoSet[tx.ID] = make(map[int]transaction.Output)
	}
	for index, output := range tx.Outputs {
		utxoSet[tx.ID][index] = output
	}
}

func removeUTXO(tx transaction.Transaction) {
	mu.Lock()
	defer mu.Unlock()
	for _, input := range tx.Inputs {
		delete(utxoSet[input.PrevTxID], input.OutputIndex)
		if len(utxoSet[input.PrevTxID]) == 0 {
			delete(utxoSet, input.PrevTxID)
		}
	}
}

func propagateTransaction(tx transaction.Transaction, peers map[int]*network.Peer) {
	for _, peer := range peers {
		peer.SendMessage(tx)
	}
}

func validateTransaction(tx transaction.Transaction) bool {
	mu.Lock()
	defer mu.Unlock()

	inputSum := 0
	outputSum := 0

	for _, input := range tx.Inputs {
		utxo, exists := utxoSet[input.PrevTxID][input.OutputIndex]
		if !exists || utxo.Value <= 0 {
			return false
		}
		inputSum += utxo.Value
	}

	for _, output := range tx.Outputs {
		outputSum += output.Value
	}

	return inputSum >= outputSum
}
