package tx

// TXInput holds the transaction output
type TXOutput struct {
	Value  int
	PubKey string
}

func (txout *TXOutput) CanBeUnlockedOutput(data string) bool {
	return txout.PubKey == data
}
