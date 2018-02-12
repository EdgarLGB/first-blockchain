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
	Data          []byte
	Hash          []byte
	PrevBlockHash []byte
	Nonce         int
}

func (b *Block) String() string {
	return fmt.Sprintf("Prev. hash: %x\n", b.PrevBlockHash) + fmt.Sprintf("Data: %s\n", b.Data) + fmt.Sprintf("Hash: %x\n", b.Hash) + fmt.Sprintf("Nonce: %d\n", b.Nonce)
}

func (b *Block) setHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{timestamp, b.Data, b.PrevBlockHash}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), []byte{}, prevBlockHash, 0}
	// do some difficult proof of work
	pow := NewProofOfWork(block)
	pow.Run()
	return block
}
