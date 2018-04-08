package main

import (
	"fmt"
	blockChain "myBlockChain/src"
	)

func main() {
	cli := blockChain.CLI{}
	//fmt.Printf("'bobo' has a balance of %d \n", cli.GetBalance("bobo"))
	cli.Send("bobo", "ying", 5)
	fmt.Printf("'bobo' has a balance of %d \n", cli.GetBalance("bobo"))
	fmt.Printf("'ying' has a balance of %d \n", cli.GetBalance("ying"))
}
