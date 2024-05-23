package concensus

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"blockchain-hello-golang/block"
)

const targetTimePerBlock = 10 * 60 // 10 minutes
const blocksPerAdjustment = 2016

func calculateHash(b block.Block) string {
	record := strconv.Itoa(b.Index) + b.PrevHash + strconv.FormatInt(b.Timestamp, 10) + b.MerkleRoot + strconv.Itoa(b.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return fmt.Sprintf("%x", hashed)
}

func solvePuzzle(b block.Block, difficulty int) int {
	target := 1 << (256 - difficulty)
	nonce := 0
	for {
		hash := calculateHash(b)
		hashInt, _ := strconv.ParseUint(hash, 16, 64)
		if hashInt < uint64(target) {
			return nonce
		}
		nonce++
		b.Nonce = nonce
	}
}

func mineBlock(transactions []block.Transaction, prevBlock block.Block, difficulty int) block.Block {
	newBlock := block.Block{
		Index:        prevBlock.Index + 1,
		PrevHash:     prevBlock.Hash,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		Difficulty:   difficulty,
	}
	newBlock.MerkleRoot = block.CalculateMerkleRoot(newBlock.Transactions)
	newBlock.Nonce = solvePuzzle(newBlock, difficulty)
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

func adjustDifficulty(chain []block.Block) int {
	if len(chain)%blocksPerAdjustment != 0 {
		return chain[len(chain)-1].Difficulty
	}

	lastAdjustmentBlock := chain[len(chain)-blocksPerAdjustment]
	expectedTime := targetTimePerBlock * blocksPerAdjustment
	actualTime := chain[len(chain)-1].Timestamp - lastAdjustmentBlock.Timestamp

	if actualTime < expectedTime/2 {
		return lastAdjustmentBlock.Difficulty + 1
	} else if actualTime > expectedTime*2 {
		return lastAdjustmentBlock.Difficulty - 1
	}
	return lastAdjustmentBlock.Difficulty
}
