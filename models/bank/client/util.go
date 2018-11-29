package client

import (
	bank "github.com/PhenixChain/phenix-go/models/bank"
	sdk "github.com/PhenixChain/phenix-go/models/types"
)

// create the sendTx msg
func CreateMsg(from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) sdk.Msg {
	input := bank.NewInput(from, coins)
	output := bank.NewOutput(to, coins)
	msg := bank.NewMsgSend([]bank.Input{input}, []bank.Output{output})
	return msg
}
