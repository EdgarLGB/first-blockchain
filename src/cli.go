package src

import "fmt"

type CLI struct {

}

func (cli *CLI) Send(from, to string, amount int) {
	bc := NewBlockChain()
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CLI) GetBalance(name string) int {
	bc := NewBlockChain()
	defer bc.db.Close()

	bal := 0
	txs := bc.FindUnspentTransactions(name)
	for _, tx := range txs {
		for _, out := range tx.Vout {
			if out.CanBeUnlocked(name) {
				bal += out.Value
			}
		}
	}
	return bal
}
