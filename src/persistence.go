package src

import (
	"bytes"
	"encoding/gob"
)

/*
Transform a block into a byte array so that it can be persisted
 */
func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(b); err != nil {
		return []byte{}, err
	} else {
		return result.Bytes(), nil
	}
}

/*
Transform a byte array into a block
 */
func DeserializeBlock(b []byte) (*Block, error) {
	block := &Block{}
	decoder := gob.NewDecoder(bytes.NewReader(b))
	if err := decoder.Decode(block); err != nil {
		return nil, err
	} else {
		return block, nil
	}
}
