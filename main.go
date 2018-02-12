package main

import (
	"fmt"
	blockChain "myBlockChain/src"
	)

func main() {
	bc := blockChain.NewBlockChain()
	//bc.AddBlock("Bobo send 1 euro to Ying.")
	//bc.AddBlock("Ying send 2 euro to Bobo.")
	i := bc.Iterator()
	for {
		block := i.Next()
		if block == nil {
			break
		}
		fmt.Println(block)
	}
}
