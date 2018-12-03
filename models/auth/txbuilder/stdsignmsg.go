package txbuilder

import (
	"github.com/PhenixChain/phenix-go/models/auth"
	sdk "github.com/PhenixChain/phenix-go/models/types"
)

// StdSignMsg is a convenience structure for passing along
// a Msg with the other requirements for a StdSignDoc before
// it is signed. For use in the CLI.
type StdSignMsg struct {
	ChainID  string      `json:"chain_id"`
	Sequence int64       `json:"sequence"`
	Fee      auth.StdFee `json:"fee"`
	Msgs     []sdk.Msg   `json:"msgs"`
	Memo     string      `json:"memo"`
}

// get message bytes
func (msg StdSignMsg) Bytes() []byte {
	return auth.StdSignBytes(msg.ChainID, msg.Sequence, msg.Fee, msg.Msgs, msg.Memo)
}
