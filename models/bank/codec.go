package bank

import (
	"github.com/PhenixChain/phenix-go/models/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/MsgSend", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
