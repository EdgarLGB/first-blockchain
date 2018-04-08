package src

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"encoding/hex"
)

const BlocksDbPath = "blocks.db"
const BlocksBucketName = "blocks-2"

type BlockChain struct {
	tip []byte	// the hash value of the last block in this chain
	db *bolt.DB  // we use bolt as database
}

/**
  Will be changed to add transactions into the blockchain
 */
func (bc *BlockChain) MineBlock(txs []*Transaction) {
	// mine a new block according to the hash value of the last block
	block := NewBlock(txs, bc.tip)
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

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
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
			coinbase := NewCoinbaseTransaction()
			genesis := NewGenesisBlock(coinbase)
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
	unspentOutputs := make(map[string][]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

	Work:
		for _, tx := range unspentTxs {
			txId := hex.EncodeToString(tx.ID)
			for outId, out := range tx.Vout {
				if out.CanBeUnlocked(address) && accumulated < amount{
					accumulated += out.Value
					unspentOutputs[txId] = append(unspentOutputs[txId], outId)
				}
				if accumulated >= amount {
					break Work
				}
			}
		}
	return accumulated, unspentOutputs
}

type BlockChainIterator struct {
	currentHash []byte
	db *bolt.DB
}

func (bc *BlockChain) Iterator() *BlockChainIterator  {
	return &BlockChainIterator{bc.tip, bc.db}
}

/**
 TODO It seems incorrect to include transactions containing the spent output. Should return directly the unspent outputs
 */
func (bc *BlockChain) FindUnspentTransactions(addr string) []*Transaction {
	var unspentTXs []*Transaction
	spentTXOs := make(map[string][]int)
	it := bc.Iterator()

	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.ID)
		Output:
			for outId, out := range tx.Vout {
				// check if already spent
				if spentTXOs[txId] != nil {
					for _, id := range spentTXOs[txId] {
						if id == outId {
							continue Output
						}
					}
				}
				if out.CanBeUnlocked(addr) {
					unspentTXs = append(unspentTXs, tx)
				}
			}
			if tx.isCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutput(addr) {
						inTxId := hex.EncodeToString(in.Txid)
						spentTXOs[inTxId] = append(spentTXOs[inTxId], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
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