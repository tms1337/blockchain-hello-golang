package block

import (
	"blockchain-hello-golang/transaction"
	"crypto/sha256"
	"fmt"
	"strconv"
)

type Block struct {
	Index        int
	PrevHash     string
	Timestamp    int64
	Transactions []transaction.Transaction
	MerkleRoot   string
	Nonce        int
	Hash         string
	Difficulty   int
}

func calculateMerkleRoot(transactions []transaction.Transaction) string {
	var txHashes []string
	for _, tx := range transactions {
		txHashes = append(txHashes, tx.ID)
	}
	return merkle(txHashes)
}

func merkle(hashes []string) string {
	if len(hashes) == 1 {
		return hashes[0]
	}

	var newLevel []string
	for i := 0; i < len(hashes); i += 2 {
		if i+1 == len(hashes) {
			newLevel = append(newLevel, hashes[i])
		} else {
			combined := hashes[i] + hashes[i+1]
			hash := sha256.Sum256([]byte(combined))
			newLevel = append(newLevel, fmt.Sprintf("%x", hash))
		}
	}
	return merkle(newLevel)
}

func createBlockHeader(index int, prevHash string, merkleRoot string, timestamp int64, difficulty int, nonce int) Block {
	return Block{
		Index:      index,
		PrevHash:   prevHash,
		Timestamp:  timestamp,
		MerkleRoot: merkleRoot,
		Nonce:      nonce,
		Difficulty: difficulty,
	}
}

func validateBlock(block, prevBlock Block) bool {
	if block.PrevHash != prevBlock.Hash {
		return false
	}
	if calculateHash(block) != block.Hash {
		return false
	}
	if calculateMerkleRoot(block.Transactions) != block.MerkleRoot {
		return false
	}
	return true
}

func calculateHash(b Block) string {
	record := strconv.Itoa(b.Index) + b.PrevHash + strconv.FormatInt(b.Timestamp, 10) + b.MerkleRoot + strconv.Itoa(b.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return fmt.Sprintf("%x", hashed)
}
