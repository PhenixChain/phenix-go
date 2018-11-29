/* package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/PhenixChain/phenix-go/models/auth"
	"github.com/PhenixChain/phenix-go/models/bank"
	"github.com/PhenixChain/phenix-go/models/codec"
	"github.com/PhenixChain/phenix-go/models/types"

	"github.com/PhenixChain/phenix-go/models/auth/txbuilder"
	"github.com/PhenixChain/phenix-go/models/bank/client"
	"github.com/PhenixChain/phenix-go/models/crypto/keys/hd"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tmlibs/bech32"
)

var Cdc = codec.New()

func init() {
	codec.RegisterCrypto(Cdc)
	types.RegisterCodec(Cdc)
	bank.RegisterCodec(Cdc)
	auth.RegisterCodec(Cdc)
}

func main() {
	GenKey()
	SendTX()
}

// GenKey ...
func GenKey() {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		log.Fatalln(err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatalln(err)
	}
	seed := bip39.NewSeed(mnemonic, "")
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, "44'/118'/0'/0/0")
	if err != nil {
		log.Fatalln(err)
	}

	prik := secp256k1.PrivKeySecp256k1(derivedPriv)
	pubk := prik.PubKey()
	Addr, err := bech32.ConvertAndEncode("adr", pubk.Address().Bytes())
	if err != nil {
		log.Fatalln(err)
	}
	PubKey, _ := bech32.ConvertAndEncode("pub", pubk.Bytes())
	fmt.Println("Address:"+Addr, "PublicKey:"+PubKey)

	//_, bz, err := bech32.DecodeAndConvert(Addr)
	//hexPubKey := append([]byte("account:"), bz...)
	//if err != nil {
	//log.Fatalln(err)
	//}
	//fmt.Println("Hex PublicKey:" + hex.EncodeToString(hexPubKey))

	fmt.Println("Mnemonic:" + mnemonic)
}

// SendTX ...
func SendTX() {
	//fromAdr := "adr12fxqmhv9steldtqykkjm2emql8eqfvw6am76xj"
	fromAdr := "adr1ttsph4qv93hllu8spl026s0rfmwhfl9d6fenyw"
	toAdr := "adr1yrd22rg0hq3wkj4jwv0s8z8xp9fpnah8dd5u59"
	from, err := types.AccAddressFromBech32(fromAdr)
	if err != nil {
		log.Fatalln(err)
	}
	to, err := types.AccAddressFromBech32(toAdr)
	if err != nil {
		log.Fatalln(err)
	}

	coins, err := types.ParseCoins("1coin1")
	if err != nil {
		log.Fatalln(err)
	}
	msg := client.CreateMsg(from, to, coins)

	tb := txbuilder.StdSignMsg{
		ChainID:       "phenix",
		AccountNumber: 2,
		Sequence:      4,
		Memo:          "",
		Msgs:          []types.Msg{msg},
		Fee:           auth.NewStdFee(200000, types.Coin{}),
	}
	sign, err := buildAndSign(tb)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(hex.EncodeToString(sign))
}

func buildAndSign(msg txbuilder.StdSignMsg) ([]byte, error) {

	//mnemonic := "bounce prevent cross remind lunch pitch project dragon firm stove labor bicycle phrase giggle cliff huge betray mask ecology gloom access alarm yellow tuna"
	mnemonic := "unfair subway explain reward shrug cement dial junk twin vital badge sing lift chair cage interest rack fault feature original acoustic vote sheriff car"
	seed := bip39.NewSeed(mnemonic, "")
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, "44'/118'/0'/0/0")
	if err != nil {
		log.Fatalln(err)
	}

	priv := secp256k1.PrivKeySecp256k1(derivedPriv)

	sigBytes, err := priv.Sign(msg.Bytes())
	if err != nil {
		log.Fatalln(err)
	}
	pubkey := priv.PubKey()

	sig := auth.StdSignature{
		AccountNumber: msg.AccountNumber,
		Sequence:      msg.Sequence,
		PubKey:        pubkey,
		Signature:     sigBytes,
	}
	return Cdc.MarshalJSON(auth.NewStdTx(msg.Msgs, msg.Fee, []auth.StdSignature{sig}, msg.Memo))
}
*/