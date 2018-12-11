package bank

import (
	"encoding/json"

	sdk "github.com/PhenixChain/phenix-go/models/types"
)

// MsgSend - high level transaction of the coin module
type MsgSend struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

var _ sdk.Msg = MsgSend{}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(in []Input, out []Output) MsgSend {
	return MsgSend{Inputs: in, Outputs: out}
}

// Implements Msg.
func (msg MsgSend) GetSignBytes() []byte {
	var inputs, outputs []json.RawMessage
	for _, input := range msg.Inputs {
		inputs = append(inputs, input.GetSignBytes())
	}
	for _, output := range msg.Outputs {
		outputs = append(outputs, output.GetSignBytes())
	}
	b, err := msgCdc.MarshalJSON(struct {
		Inputs  []json.RawMessage `json:"inputs"`
		Outputs []json.RawMessage `json:"outputs"`
	}{
		Inputs:  inputs,
		Outputs: outputs,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

//----------------------------------------
// Input

// Transaction Input
type Input struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// Return bytes to sign for Input
func (in Input) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(in)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bin)
}

// NewInput - create a transaction input, used with MsgSend
func NewInput(addr sdk.AccAddress, coins sdk.Coins) Input {
	input := Input{
		Address: addr,
		Coins:   coins,
	}
	return input
}

//----------------------------------------
// Output

// Transaction Output
type Output struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// Return bytes to sign for Output
func (out Output) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(out)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bin)
}

// NewOutput - create a transaction output, used with MsgSend
func NewOutput(addr sdk.AccAddress, coins sdk.Coins) Output {
	output := Output{
		Address: addr,
		Coins:   coins,
	}
	return output
}
