package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"

	"fmt"
)

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
	if err != nil {
		return nil, nil, err
	}
	return r, s, nil
}

func VerifySignature(pubKey *ecdsa.PublicKey, msg []byte, r, s *big.Int) bool {
	hash := sha256.Sum256(msg)
	return ecdsa.Verify(pubKey, hash[:], r, s)
}

func PublicKeyToAddress(pubKey *ecdsa.PublicKey) string {
	pubKeyBytes := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)
	sha256Hash := sha256.Sum256(pubKeyBytes)
	ripemd160Hasher := ripemd160.New()
	_, err := ripemd160Hasher.Write(sha256Hash[:])
	if err != nil {
		log.Fatal(err)
	}
	ripemd160Hash := ripemd160Hasher.Sum(nil)
	return fmt.Sprintf("%x", ripemd160Hash)
}

func Hash160(data []byte) []byte {
	sha256Hash := sha256.Sum256(data)
	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(sha256Hash[:])
	return ripemd160Hasher.Sum(nil)
}

func EncodeBase58Check(data []byte) string {
	checksum := sha256.Sum256(data)
	checksum = sha256.Sum256(checksum[:])
	fullData := append(data, checksum[:4]...)
	return base58Encode(fullData)
}

func base58Encode(data []byte) string {
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	var result []byte
	x := new(big.Int).SetBytes(data)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	for x.Cmp(zero) != 0 {
		mod := new(big.Int)
		x.DivMod(x, base, mod)
		result = append(result, alphabet[mod.Int64()])
	}
	reverse(result)
	return string(result)
}

func reverse(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
