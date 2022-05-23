package model

import (
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"

	. "github.com/rrossilli/glow/util"
)

// Account struct as it typically appears in a flow.json
type Account struct {
	Address string `json:"address"`
	PrivKey string `json:"key"`
}

// New Account
func NewAccount(address, privKey string) Account {
	return Account{
		Address: address,
		PrivKey: privKey,
	}
}

// New Account without a private key
func NewUnqualifiedAccount(address string) Account {
	return Account{
		Address: address,
	}
}

// "flow" address
func (a Account) FlowAddress() flow.Address {
	return flow.HexToAddress(a.Address)
}

// "cadence" address
func (a Account) CadenceAddress() cadence.Address {
	return cadence.Address(a.FlowAddress())
}

// Crypto private key
func (a Account) CryptoPrivateKey() crypto.PrivateKey {
	key, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, RemoveHexPrefix(a.PrivKey))
	if err != nil {
		panic(err)
	}
	return key
}

// Crypto public key
func (a Account) CryptoPublicKey() crypto.PublicKey {
	return a.CryptoPrivateKey().PublicKey()
}

// Flow addresses
func FlowAddressesFromAccounts(as []Account) []flow.Address {
	var addrs []flow.Address
	for _, v := range as {
		addrs = append(addrs, v.FlowAddress())
	}
	return addrs
}

// Cadence addresses
func CadenceAddressesFromAccounts(as []Account) []cadence.Address {
	var addrs []cadence.Address
	for _, v := range as {
		addrs = append(addrs, v.CadenceAddress())
	}
	return addrs
}
