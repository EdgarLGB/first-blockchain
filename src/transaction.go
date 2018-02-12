package src

type TXInput struct {
	
}

type TXOutput struct {
	Value int
	ScriptPubKey string
}


type Transaction struct {
	ID []byte
	Vin []TXInput
	Vout []TXOutput
}

func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	bc.FindSpendableOutputs(from, amount)
}
