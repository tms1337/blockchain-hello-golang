package transaction

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type Input struct {
	PrevTxID    string
	OutputIndex int
	ScriptSig   string
}

type Output struct {
	Value        int
	ScriptPubKey string
}

type Transaction struct {
	ID      string
	Inputs  []Input
	Outputs []Output
}

func createTransaction(inputs []Input, outputs []Output) Transaction {
	tx := Transaction{Inputs: inputs, Outputs: outputs}
	tx.ID = calculateTransactionID(tx)
	return tx
}

func validateTransaction(tx Transaction, utxoSet map[string]Output) bool {
	inputSum := 0
	outputSum := 0

	for _, input := range tx.Inputs {
		utxo, exists := utxoSet[input.PrevTxID]
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

func calculateTransactionID(tx Transaction) string {
	record := ""
	for _, input := range tx.Inputs {
		record += input.PrevTxID + string(input.OutputIndex) + input.ScriptSig
	}
	for _, output := range tx.Outputs {
		record += string(output.Value) + output.ScriptPubKey
	}
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return fmt.Sprintf("%x", hashed)
}

func validateScript(scriptSig, scriptPubKey string) bool {
	// Simplified script validation for now
	return scriptSig == scriptPubKey
}

func GenerateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

func SignMessage(privKey *ecdsa.PrivateKey, msg []byte) (r, s *big.Int, err error) {
	hash := sha256.Sum256(msg)
	r, s, err = ecdsa.Sign(rand.Reader, privKey, hash[:])
	return
}

func VerifySignature(pubKey *ecdsa.PublicKey, msg []byte, r, s *big.Int) bool {
	hash := sha256.Sum256(msg)
	return ecdsa.Verify(pubKey, hash[:], r, s)
}
