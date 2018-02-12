package src

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
)

const BlocksDbPath = "blocks.db"
const BlocksBucketName = "blocks"

type BlockChain struct {
	tip []byte	// the hash value of the last block in this chain
	db *bolt.DB
}

func (bc *BlockChain) AddBlock(data string) {
	// mine a new block according to the hash value of the last block
	block := NewBlock(data, bc.tip)
	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte([]byte(BlocksBucketName)))
		blockSerialized, err := block.Serialize()
		if err != nil {
			return fmt.Errorf("serialize a block: %s", err)
		}
		err = b.Put([]byte(block.Hash), blockSerialized)
		if err != nil {
			return fmt.Errorf("add a block in bucket: %s", err)
		}
		err = b.Put([]byte("l"), block.Hash)
		if err != nil {
			return fmt.Errorf("update the last block hash in bucket: %s", err)
		}
		bc.tip = block.Hash
		return nil
	})
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

/*
Create a new block chain and persists it into a db
 */
func NewBlockChain() *BlockChain {
	var tip []byte
	db, err := bolt.Open(BlocksDbPath, 0600, nil)
	if err != nil {
		// error needs to be bubbled up
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucketName))
		if b == nil {
			b, err := tx.CreateBucket([]byte(BlocksBucketName))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			genesis := NewGenesisBlock()
			genesisBytes, err := genesis.Serialize()
			if err != nil {
				return fmt.Errorf("serialize genesis: %s", err)
			}
			err = b.Put(genesis.Hash, genesisBytes)
			if err != nil {
				return fmt.Errorf("save block: %s", err)
			}
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				return fmt.Errorf("save block: %s", err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	return &BlockChain{tip, db}
}

func (bc *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	
}

type BlockChainIterator struct {
	currentHash []byte
	db *bolt.DB
}

func (bc *BlockChain) Iterator() *BlockChainIterator  {
	return &BlockChainIterator{bc.tip, bc.db}
}

func (bci *BlockChainIterator) Next() *Block {
	var result *Block
	bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucketName))
		blockHash := b.Get([]byte(bci.currentHash))
		block, err := DeserializeBlock(blockHash)
		result = block
		if err != nil {
			return fmt.Errorf("deserialize block error: %s", err)
		}
		bci.currentHash = block.PrevBlockHash
		return nil
	})
	return result
}