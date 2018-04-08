package src

import (
	"math/big"
	"fmt"
	"crypto/sha256"
	"bytes"
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1) // The current target is based on 1
	target.Lsh(target, uint(256-TargetBits))
	return &ProofOfWork{block, target}
}

/*
Join the block header with the nonce
 */
func (pow *ProofOfWork) prepareData(nonce int) ([]byte) {
	var tb []byte
	for _, t := range pow.block.Transactions {
		tb = append(tb, t.serialize()...)
	}
	// the int64 needs to be converted to hex and then to be concatenated
	return bytes.Join([][]byte{
		tb,
		pow.block.PrevBlockHash,
		IntToHex(pow.block.Timestamp),
		IntToHex(int64(TargetBits)),
		IntToHex(int64(nonce))}, []byte{})
}

/**
  Do the POW to mine the block
 */
func (pow *ProofOfWork) Run() {
	var hashInt big.Int
	var hash [32]byte // the type of hash value is defined by result of the sha256 function
	nonce := 0

	fmt.Printf("Start mining the block \"%s\"\n", pow.block.Hash)
	for nonce < MaxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			// the nonce found
			fmt.Printf("\r%x\n\n", hash)
			break
		} else {
			nonce++
		}
	}
	pow.block.Hash = hash[:]
	pow.block.Nonce = nonce
}

func (pow *ProofOfWork) Validate() bool {
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	var hashInt big.Int
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}
