package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PhenixChain/phenix-go/models/auth"
	"github.com/PhenixChain/phenix-go/models/bank"
	"github.com/PhenixChain/phenix-go/models/types"

	"github.com/PhenixChain/phenix-go/models/auth/txbuilder"
	"github.com/PhenixChain/phenix-go/models/crypto/hd"
	bip39 "github.com/cosmos/go-bip39"
	amino "github.com/tendermint/go-amino"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/bech32"
)

var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
	cdc.RegisterInterface((*types.Msg)(nil), nil)
	cdc.RegisterConcrete(bank.MsgSend{}, "cosmos-sdk/MsgSend", nil)
	cdc.RegisterConcrete(auth.StdTx{}, "auth/StdTx", nil)
}

func main() {
	startTime := time.Now()

	//##################### 生成公私钥 #################################################

	//genKey()

	//#####################  发送交易  #################################################

	// 转出地址
	fromAdr := "adr12fxqmhv9steldtqykkjm2emql8eqfvw6am76xj"
	//fromAdr := "adr1ttsph4qv93hllu8spl026s0rfmwhfl9d6fenyw"

	// 转入地址
	toAdr := "adr1yrd22rg0hq3wkj4jwv0s8z8xp9fpnah8dd5u59"

	// 金额币种
	coin := "6coin1"

	// 转出地址对应的助记词
	mnemonic := "bounce prevent cross remind lunch pitch project dragon firm stove labor bicycle phrase giggle cliff huge betray mask ecology gloom access alarm yellow tuna"
	//mnemonic := "unfair subway explain reward shrug cement dial junk twin vital badge sing lift chair cage interest rack fault feature original acoustic vote sheriff car"

	// 交易序号
	sequence := uint64(0)
	if seq := getAccount(fromAdr, false); seq != "" {
		sequence, _ = strconv.ParseUint(seq, 10, 64)
	}

	sendTX(fromAdr, toAdr, coin, mnemonic, sequence)

	//#####################  查询地址余额  #################################################

	//addr := "adr12fxqmhv9steldtqykkjm2emql8eqfvw6am76xj"

	//getAccount(addr, true)

	//######################   查询交易   #################################################

	//tx := "23ECD5198818D03B6568257F1BA29BD4DC0668CAD95CD6472C14C7E09AF6AAD3"

	// 查询交易(txhash、height)
	//getTX(tx)

	//####################### 查询地址的txhash ##############################################

	//addr = "adr12fxqmhv9steldtqykkjm2emql8eqfvw6am76xj"

	//getTXByAddr(addr)

	elapsed := time.Since(startTime)
	fmt.Println("elapsed cost: ", elapsed)
}

func genKey() {
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

	fmt.Println("Address:   " + Addr)
	fmt.Println("PublicKey: " + PubKey)
	fmt.Println("Mnemonic:  " + mnemonic)
}

func getAccount(addr string, show bool) string {
	_, bz, err := bech32.DecodeAndConvert(addr)
	//hexPubKey := append([]byte("account:"), bz...) //v1.1
	hexPubKey := append([]byte{0x01}, bz...) //v1.2+
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println("Hex PublicKey:" + hex.EncodeToString(hexPubKey))

	url := `http://127.0.0.1:26657/abci_query?path="/store/acc/key"&data=0x` + hex.EncodeToString(hexPubKey)
	res := httpGet(url)

	accRes := AccountResponse{}
	err = json.Unmarshal(res, &accRes)
	if err != nil {
		log.Fatalln(err)
	}
	if accRes.Result.Response.Value != "" {
		br, err := base64.StdEncoding.DecodeString(accRes.Result.Response.Value)
		if err != nil {
			log.Fatalln(err)
		}

		accInfo := AccountInfo{}
		err = json.Unmarshal(br, &accInfo)
		if err != nil {
			log.Fatalln(err)
		}

		if show {
			fmt.Println(string(br))
		}

		return accInfo.Value.Sequence
	}

	return ""
}

func getTXByAddr(addr string) {
	_, bz, err := bech32.DecodeAndConvert(addr)
	hexPubKey := append([]byte{0x01}, bz...) //v1.2+
	if err != nil {
		log.Fatalln(err)
	}

	url := `http://127.0.0.1:26657/abci_query?path="/store/address/key"&data=0x` + hex.EncodeToString(hexPubKey)
	res := httpGet(url)

	accRes := AccountResponse{}
	err = json.Unmarshal(res, &accRes)
	if err != nil {
		log.Fatalln(err)
	}

	br, err := base64.StdEncoding.DecodeString(accRes.Result.Response.Value)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(br))
}

func getTX(tx string) {
	url := `http://127.0.0.1:26657/tx?hash=0x` + tx
	res := httpGet(url)

	tranRes := TranResponse{}
	err := json.Unmarshal(res, &tranRes)
	if err != nil {
		log.Fatalln(err)
	}

	br, err := base64.StdEncoding.DecodeString(tranRes.Result.Tx)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(br))
}

func sendTX(fromAdr, toAdr, coin, mnemonic string, sequence uint64) {
	from, err := types.AccAddressFromBech32(fromAdr)
	if err != nil {
		log.Fatalln(err)
	}
	to, err := types.AccAddressFromBech32(toAdr)
	if err != nil {
		log.Fatalln(err)
	}
	coins, err := types.ParseCoins(coin)
	if err != nil {
		log.Fatalln(err)
	}
	msg := bank.NewMsgSend(from, to, coins)

	//<!-- < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < << < < ☺
	//v                  ✰  交易费处理 ✰
	//v
	//v    	交易费上限 fees = ceil(gas * gasPrices)
	//v		如：gas Limit = 200000; gas Prices = 5Micro
	//v		fees = 200000*5*(-10^6)
	//v
	//v		gas是衡量交易需要消耗多少资源的单位
	//v		gas上限用--gas指定(推荐将gas上限设置为200000)。gas上限太小时，不够交易需要的gas
	//v		实际消耗多少gas就会花多少相应的交易费，剩余的交易费会被退还
	//v
	//v		1 coin1 = 10^3 coin1-milli
	//v		Milli = "milli"
	//v
	//v		1 coin1 = 10^6 coin1-micro
	//v		Micro = "micro"
	//v
	//v		1 coin1 = 10^9 coin1-nano
	//v		Nano = "nano"
	//v
	//v		1 coin1 = 10^12 coin1-pico
	//v		Pico = "pico"
	//v
	//v		1 coin1 = 10^15 coin1-femto
	//v		Femto = "femto"
	//v
	//v		1 coin1 = 10^18 coin1-atto
	//v		Atto = "atto"
	//v
	//☺ > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > >  -->

	fees, err := types.ParseCoins("1coin1")
	if err != nil {
		log.Fatalln(err)
	}

	tb := txbuilder.StdSignMsg{
		ChainID:  "phenix",
		Sequence: sequence,
		Memo:     "",
		Msgs:     []types.Msg{msg},
		Fee:      auth.NewStdFee(200000, fees),
	}

	sign, err := buildAndSign(tb, mnemonic)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println(hex.EncodeToString(sign))

	//Commit (Waiting for a new block)
	//url := "http://127.0.0.1:26657/broadcast_tx_commit?tx=0x" + hex.EncodeToString(sign)

	//Propose (Waiting for the proposal result)
	url := "http://127.0.0.1:26657/broadcast_tx_sync?tx=0x" + hex.EncodeToString(sign)
	fmt.Println(string(httpGet(url)))
}

func buildAndSign(msg txbuilder.StdSignMsg, mnemonic string) ([]byte, error) {
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
		PubKey:    pubkey,
		Signature: sigBytes,
	}

	tx := auth.StdTx{
		Msgs:       msg.Msgs,
		Fee:        msg.Fee,
		Signatures: []auth.StdSignature{sig},
		Memo:       msg.Memo,
	}

	return cdc.MarshalJSON(tx)
}

func httpGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

type TranResponse struct {
	Result TxResult `json:"result"`
}

type TxResult struct {
	Tx string `json:"tx"`
}

type AccountResponse struct {
	Result Result
}

type Result struct {
	Response Response
}

type Response struct {
	Value string `json:"value"`
}

type AccountInfo struct {
	Value Value
}

type Value struct {
	Address  string `json:"address"`
	Coins    []Coins
	Sequence string `json:"sequence"`
}

type Coins struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
