package blockchain

type TXInput struct {
	ID        []byte
	Output    int
	Signature string
}

type TXOutput struct {
	Value  int
	PubKey string
}

// CanUnlock checks whether the address is the owner of the input.
func (in *TXInput) CanUnlock(data string) bool {
	return in.Signature == data
}

// CanBeUnlocked checks whether the address is the owner of the output.
func (out *TXOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
