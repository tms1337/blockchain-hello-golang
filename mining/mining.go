package mining

import (
	"sort"

	"blockchain-hello-golang/block"
	"blockchain-hello-golang/transaction"
)

const initialReward = 50
const halvingInterval = 210000

func calculateReward(height int) int {
	halvings := height / halvingInterval
	return initialReward >> halvings
}

func createCoinbaseTx(minerAddress string, reward int) transaction.Transaction {
	return transaction.Transaction{
		ID:     "coinbase",
		Inputs: []transaction.Input{},
		Outputs: []transaction.Output{
			{Value: reward, ScriptPubKey: minerAddress},
		},
	}
}

func selectTransactions(mempool []transaction.Transaction, blockSize int) []transaction.Transaction {
	var selectedTxs []transaction.Transaction
	currentSize := 0

	// Sort transactions by fee in descending order
	sort.Slice(mempool, func(i, j int) bool {
		return calculateTxFee(mempool[i]) > calculateTxFee(mempool[j])
	})

	for _, tx := range mempool {
		txSize := calculateTxSize(tx)
		if currentSize+txSize > blockSize {
			break
		}
		selectedTxs = append(selectedTxs, tx)
		currentSize += txSize
	}
	return selectedTxs
}

func calculateTxSize(tx transaction.Transaction) int {
	// Simplified size calculation
	return len(tx.ID) + len(tx.Inputs)*100 + len(tx.Outputs)*100
}

func calculateTxFee(tx transaction.Transaction) int {
	inputSum := 0
	outputSum := 0

	for _, input := range tx.Inputs {
		// Assuming utxoSet is globally accessible
		utxo := utxoSet[input.PrevTxID][input.OutputIndex]
		inputSum += utxo.Value
	}

	for _, output := range tx.Outputs {
		outputSum += output.Value
	}

	return inputSum - outputSum
}

func includeCoinbaseTx(block *block.Block, minerAddress string, reward int) {
	coinbaseTx := createCoinbaseTx(minerAddress, reward)
	block.Transactions = append([]transaction.Transaction{coinbaseTx}, block.Transactions...)
}
