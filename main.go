package main

import (
	"fmt"
	blockChain "myBlockChain/src"
	)

func main() {
	bc := blockChain.NewBlockChain()
	bc.AddBlock("Bobo send 1 euro to Ying.")
	bc.AddBlock("Ying send 2 euro to Bobo.")
	for _, block := range bc.GetBlocks() {
		fmt.Println(block)
	}
}
