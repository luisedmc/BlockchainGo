package main

// TXInput holds the transaction output
type TXOutput struct {
	Amount int
	PubKey string
}

func (txout *TXOutput) CanUnlockOutput(data string) bool {
	return txout.PubKey == data
}
