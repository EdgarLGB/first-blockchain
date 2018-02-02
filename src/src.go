package src

import (
	"strconv"
	"bytes"
	"crypto/sha256"
	"time"
	"fmt"
	"math/big"
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

type BlockChain struct {
	blocks []*Block
}

func (bc *BlockChain) GetBlocks() []*Block {
	return bc.blocks
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	block := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, block)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockChain() *BlockChain {
	block := NewGenesisBlock()
	return &BlockChain{[]*Block{block}}
}

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
	// the int64 needs to be converted to hex and then to be concatenated
	return bytes.Join([][]byte{
		pow.block.Data,
		pow.block.PrevBlockHash,
		IntToHex(pow.block.Timestamp),
		IntToHex(int64(TargetBits)),
		IntToHex(int64(nonce))}, []byte{})
}

func (pow *ProofOfWork) Run() {
	var hashInt big.Int
	var hash [32]byte // the type of hash value is defined by result of the sha256 function
	nonce := 0

	fmt.Printf("Start mining the block \"%s\"\n", pow.block.Data)
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
