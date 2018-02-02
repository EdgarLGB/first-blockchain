package src

import (
	"bytes"
	"encoding/gob"
)

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(b); err != nil {
		return []byte{}, err
	} else {
		return result.Bytes(), nil
	}
}

func DeserializeBlock(b []byte) (*Block, error) {
	var block *Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	if err := decoder.Decode(block); err != nil {
		return nil, err
	} else {
		return block, nil
	}
}
