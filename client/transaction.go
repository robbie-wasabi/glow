package client

import (
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"

	. "github.com/rrossilli/glow/model"
)

type Tx struct {
	FlowTransaction flow.Transaction
	Client          GlowClient
}

// Create new Transaction
func (c GlowClient) tx(
	cdc []byte,
	proposer Account,
	payer Account,
	authorizers []Account,
	args ...cadence.Value,
) (*Tx, error) {
	tx, err := c.Services.Transactions.Build(
		proposer.FlowAddress(),
		FlowAddressesFromAccounts(authorizers),
		proposer.FlowAddress(),
		0, // todo: which key?
		cdc,
		"", // is this important?
		c.gasLimit,
		args,
		c.network,
		true,
	)
	if err != nil {
		return nil, err
	}

	wrappedTx := Tx{
		FlowTransaction: *tx.FlowTransaction(),
		Client:          c,
	}

	return &wrappedTx, nil
}

// Create new Transaction from cadence string
func (c GlowClient) NewTx(
	cdc string,
	proposer Account,
	payer Account,
	authorizers []Account,
	args []cadence.Value,
) (*Tx, error) {
	tx, err := c.tx(
		[]byte(cdc),
		proposer,
		payer,
		authorizers,
		args...,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Create new Transaction from file path
func (c GlowClient) NewTxFromFile(
	file string,
	proposer Account,
	payer Account,
	authorizers []Account,
	args ...cadence.Value,
) (*Tx, error) {
	cdc, err := c.CadenceFromFile(file)
	if err != nil {
		return nil, err
	}

	tx, err := c.tx(
		[]byte(cdc),
		proposer,
		payer,
		authorizers,
		args...,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Create new crypto signer
func (c GlowClient) newInMemorySigner(privKey string) (crypto.Signer, error) {
	pk, err := c.NewPrivateKeyFromHex(privKey)
	if err != nil {
		return nil, err
	}

	return crypto.NewInMemorySigner(pk, c.HashAlgo), nil
}

// Sign a Transaction with multiple signers
func (tx *Tx) Sign(authorizers ...Account) error {
	// map to slice of crypto signers
	var signers []crypto.Signer
	for _, a := range authorizers {
		s, err := tx.Client.newInMemorySigner(a.PrivKey)
		if err != nil {
			return err
		}
		signers = append(signers, s)
	}

	// map to slice of flow addresses
	var addresses []flow.Address
	for _, a := range authorizers {
		addresses = append(addresses, a.FlowAddress())
	}

	// sign transaction with each signer
	for i := len(addresses) - 1; i >= 0; i-- {
		signerAddress := addresses[i]
		signer := signers[i]

		if i == 0 {
			err := tx.FlowTransaction.SignEnvelope(signerAddress, 0, signer)
			if err != nil {
				return err
			}
		} else {
			err := tx.FlowTransaction.SignPayload(signerAddress, 0, signer)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Send a signed Transaction
func (signedTx Tx) Send() (*flow.TransactionResult, error) {
	txBytes := []byte(fmt.Sprintf("%x", signedTx.FlowTransaction.Encode()))
	_, res, err := signedTx.Client.Services.Transactions.SendSigned(txBytes, true)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}

	return res, nil
}

// Sign and send a Transaction from string. Assumes that proposer is also payer
func (c GlowClient) SignAndSendTx(
	cdc string,
	proposer Account,
	args ...cadence.Value,
) (*flow.TransactionResult, error) {
	tx, err := c.NewTx(
		cdc,
		proposer,
		proposer,
		[]Account{proposer},
		args,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Sign(proposer); err != nil {
		return nil, err
	}

	res, err := tx.Send()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Sign and send a Transaction at a specified file path. Assumes that proposer is also payer
func (c GlowClient) SignAndSendTxFromFile(
	file string,
	proposer Account,
	args ...cadence.Value,
) (*flow.TransactionResult, error) {
	tx, err := c.NewTxFromFile(
		file,
		proposer,
		proposer,
		[]Account{proposer},
		args...,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Sign(proposer); err != nil {
		return nil, err
	}

	res, err := tx.Send()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Sign and send a Transaction at
func (c GlowClient) SignMultiAndSendTx(
	cdc string,
	proposer Account,
	payer Account,
	authorizers []Account,
	args ...cadence.Value,
) (*flow.TransactionResult, error) {
	tx, err := c.NewTx(
		cdc,
		proposer,
		payer,
		authorizers,
		args,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Sign(authorizers...); err != nil {
		return nil, err
	}

	res, err := tx.Send()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Sign and send a Transaction at a specified file path.
func (c GlowClient) SignMultiAndSendTxFromFile(
	file string,
	proposer Account,
	payer Account,
	authorizers []Account,
	args ...cadence.Value,
) (*flow.TransactionResult, error) {
	tx, err := c.NewTxFromFile(
		file,
		proposer,
		payer,
		authorizers,
		args...,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Sign(authorizers...); err != nil {
		return nil, err
	}

	res, err := tx.Send()
	if err != nil {
		return nil, err
	}

	return res, nil
}
