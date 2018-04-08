package src

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"encoding/hex"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

const initialValue = 10
const to = "bobo"
const scriptSignature = "bobo-private-key"

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

func (txi *TXInput) CanUnlockOutput(addr string) bool {
	return addr == txi.ScriptSig
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

func (txo *TXOutput) CanBeUnlocked(addr string) bool {
	return addr == txo.ScriptPubKey
}

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx *Transaction) setId() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (t *Transaction) serialize() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, t)
	return buf.Bytes()
}

func (t *Transaction) isCoinbase() bool {
	return len(t.Vin) == 1 && t.Vin[0].Vout == -1 && len(t.Vout) == 1
}

func NewCoinbaseTransaction() *Transaction {
	txin := TXInput{[]byte{}, -1, scriptSignature}
	txout := TXOutput{initialValue, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	return &tx
}

func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	balance, unspentOutput := bc.FindSpendableOutputs(from, amount)
	if balance < amount {
		log.Panic(fmt.Sprintf("%s has only %d, not enough balance for this transaction", from, balance))
	}

	// build inputs
	for txid, outputs := range unspentOutput {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			fmt.Errorf("error when decoding %d", txid)
		}
		for _, outId := range outputs {
			inputs = append(inputs, TXInput{txID, outId, from})
		}
	}

	// build outputs
	outputs = append(outputs, TXOutput{amount, to})
	if amount < balance {
		outputs = append(outputs, TXOutput{balance - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.setId()

	return &tx
}
