package tx

// TXInput holds the transaction input
type TXInput struct {
	ID        []byte
	Output    int
	Signature string
}

func (txin *TXInput) CanBeUnlockedInput(data string) bool {
	return txin.Signature == data
}
