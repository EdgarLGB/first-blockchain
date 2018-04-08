package src

import (
	"strconv"
	"bytes"
	"crypto/sha256"
	"time"
	"fmt"
	"math"
)

const TargetBits = 24
const MaxNonce = math.MaxInt64

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction	// should be changed to Transaction
	Hash          []byte
	PrevBlockHash []byte
	Nonce         int
}

func (b *Block) String() string {
	return fmt.Sprintf("Prev. hash: %x\n", b.PrevBlockHash) + fmt.Sprintf("Transactions: %s\n", b.Transactions) + fmt.Sprintf("Hash: %x\n", b.Hash) + fmt.Sprintf("Nonce: %d\n", b.Nonce)
}

func (b *Block) setHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	var tb []byte
	for _, t := range b.Transactions  {
		tb = append(tb, t.serialize()...)
	}
	headers := bytes.Join([][]byte{timestamp, tb, b.PrevBlockHash}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(ts []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), ts, []byte{}, prevBlockHash, 0}
	// do some difficult proof of work
	pow := NewProofOfWork(block)
	pow.Run()
	return block
}
