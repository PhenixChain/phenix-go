package types

// Transactions messages must fulfill the Msg
type Msg interface {

	// Get the canonical byte representation of the Msg.
	GetSignBytes() []byte
}

//__________________________________________________________

// Transactions objects must fulfill the Tx
type Tx interface {

	// Gets the Msg.
	GetMsgs() []Msg
}

//__________________________________________________________

// TxDecoder unmarshals transaction bytes
type TxDecoder func(txBytes []byte) (Tx, Error)
