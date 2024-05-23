package fork

import (
	"blockchain-hello-golang/block"
	"blockchain-hello-golang/transaction"
)

func selectLongestChain(chains [][]block.Block) []block.Block {
	var longestChain []block.Block
	maxPoW := 0
	for _, chain := range chains {
		pow := calculatePoW(chain)
		if pow > maxPoW {
			longestChain = chain
			maxPoW = pow
		}
	}
	return longestChain
}

func calculatePoW(chain []block.Block) int {
	pow := 0
	for _, blk := range chain {
		pow += blk.Difficulty
	}
	return pow
}

func handleOrphanedBlocks(orphanedBlocks []block.Block, chain []block.Block, utxoSet map[string]map[int]transaction.Output) []block.Block {
	for _, blk := range orphanedBlocks {
		if isValidBlock(blk, chain) {
			chain = append(chain, blk)
			updateUTXOSet(blk, utxoSet)
		}
	}
	return chain
}

func isValidBlock(blk block.Block, chain []block.Block) bool {
	if len(chain) == 0 {
		return blk.Index == 0
	}
	prevBlock := chain[len(chain)-1]
	return blk.PrevHash == prevBlock.Hash && block.calculateHash(blk) == blk.Hash
}

func updateUTXOSet(blk block.Block, utxoSet map[string]map[int]transaction.Output) {
	for _, tx := range blk.Transactions {
		for index, output := range tx.Outputs {
			if _, exists := utxoSet[tx.ID]; !exists {
				utxoSet[tx.ID] = make(map[int]transaction.Output)
			}
			utxoSet[tx.ID][index] = output
		}
		for _, input := range tx.Inputs {
			delete(utxoSet[input.PrevTxID], input.OutputIndex)
			if len(utxoSet[input.PrevTxID]) == 0 {
				delete(utxoSet, input.PrevTxID)
			}
		}
	}
}
