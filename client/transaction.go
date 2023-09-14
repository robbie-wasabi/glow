package client

import (
	"context"
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-cli/flowkit"
	"github.com/onflow/flow-cli/flowkit/transactions"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"

	. "github.com/rrossilli/glow/model"
)

// Tx struct encapsulates the logic to build, sign and send Flow transactions.
type Tx struct {
	ctx         context.Context
	script      flowkit.Script // FlowKit's "script" type handles reads from the chain.
	payer       Account
	proposer    Account
	authorizers []Account
	client      *GlowClient
}

// newTx is a private helper to deduplicate logic in public constructors.
func (c *GlowClient) newTx(
	code []byte,
	proposer Account,
	args ...cadence.Value,
) *Tx {
	return &Tx{
		script: flowkit.Script{
			Code: code,
			Args: args,
		},
		proposer: proposer,
		payer:    proposer,
		authorizers: []Account{
			proposer,
		},
		client: c,
	}
}

// NewTx creates a new unsigned transaction using byte slice.
// Assumes proposer is also the gas payer and sole authorizer.
func (c *GlowClient) NewTx(
	cdc []byte,
	proposer Account,
	args ...cadence.Value,
) *Tx {
	return c.newTx(cdc, proposer, args...)
}

// NewTxFromString creates a new unsigned transaction using string code.
// Assumes proposer is also the gas payer and sole authorizer.
func (c *GlowClient) NewTxFromString(
	cdc string,
	proposer Account,
	args ...cadence.Value,
) *Tx {
	return c.newTx([]byte(c.replaceImportAddresses(cdc)), proposer, args...)
}

// NewTxFromFile creates a new unsigned transaction from a file.
// Assumes proposer is also the gas payer and sole authorizer.
func (c *GlowClient) NewTxFromFile(
	file string,
	proposer Account,
	args ...cadence.Value,
) (*Tx, error) {
	cdc, err := c.CadenceFromFile(file)
	if err != nil {
		return nil, fmt.Errorf("tx not found at: %s", file)
	}
	return c.newTx([]byte(cdc), proposer, args...), nil
}

// WithContext adds context to a transaction.
func (t *Tx) WithContext(ctx context.Context) *Tx {
	t.ctx = ctx
	return t
}

// Args specifies arguments for a transaction.
func (t *Tx) Args(args ...cadence.Value) *Tx {
	t.script.Args = args
	return t
}

// AddArg adds a single argument to the transaction.
func (t *Tx) AddArg(arg cadence.Value) *Tx {
	t.script.Args = append(t.script.Args, arg)
	return t
}

// Payer specifies who pays for the transaction.
func (t *Tx) Payer(p Account) *Tx {
	t.payer = p
	t.authorizers = append(t.authorizers, p)
	return t
}

// Proposer specifies who proposes the transaction.
func (t *Tx) Proposer(p Account) *Tx {
	t.proposer = p
	t.authorizers = append(t.authorizers, p)
	return t
}

// Authorizers sets the authorizers of the transaction.
func (t *Tx) Authorizers(a ...Account) *Tx {
	t.authorizers = a
	return t
}

// AddAuthorizer appends an authorizer to the transaction.
func (t *Tx) AddAuthorizer(a Account) *Tx {
	t.authorizers = append(t.authorizers, a)
	return t
}

type SignedTx struct {
	ctx    context.Context
	flowTx *transactions.Transaction
	client *GlowClient
}

// Create new crypto signer
func (c *GlowClient) newInMemorySigner(privKey string) (crypto.Signer, error) {
	pk, err := c.NewPrivateKeyFromHex(privKey)
	if err != nil {
		return nil, err
	}

	signer, err := crypto.NewInMemorySigner(pk, c.HashAlgo)
	if err != nil {
		return nil, err
	}

	return signer, nil
}

// Sign tx with key at index 0. Use SignTxWithKey to specify key index
func (t *Tx) Sign() (*SignedTx, error) {

	// map to slice of crypto signers
	var signers []crypto.Signer
	for _, a := range t.authorizers {
		s, err := t.client.newInMemorySigner(a.PrivKey)
		if err != nil {
			return nil, err
		}
		signers = append(signers, s)
	}

	// map to slice of flow addresses
	var addresses []flow.Address
	for _, a := range t.authorizers {
		addresses = append(addresses, a.FlowAddress())
	}

	var txAddresses = transactions.AddressesRoles{
		Proposer:    t.proposer.FlowAddress(),
		Payer:       t.payer.FlowAddress(),
		Authorizers: FlowAddressesFromAccounts(t.authorizers),
	}

	flowTx, err := t.client.FlowKit.BuildTransaction(t.ctx,
		txAddresses,
		0,
		t.script,
		t.client.gasLimit,
	)
	if err != nil {
		return nil, err
	}

	// sign transaction with each signer
	for i := len(addresses) - 1; i >= 0; i-- {
		signerAddress := addresses[i]
		signer := signers[i]

		if i == 0 {
			err := flowTx.FlowTransaction().SignEnvelope(signerAddress, 0, signer)
			if err != nil {
				return nil, err
			}
		} else {
			err := flowTx.FlowTransaction().SignPayload(signerAddress, 0, signer)
			if err != nil {
				return nil, err
			}
		}
	}

	return &SignedTx{
		ctx:    t.ctx,
		flowTx: flowTx,
		client: t.client,
	}, err
}

// Send a signed Transaction
func (signedTx *SignedTx) Send() (*flow.TransactionResult, error) {
	_, res, err := signedTx.client.FlowKit.SendSignedTransaction(signedTx.ctx, signedTx.flowTx)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}

	return res, nil
}

// Sign and send a transaction
func (tx *Tx) SignAndSend() (*flow.TransactionResult, error) {
	signedTx, err := tx.Sign()
	if err != nil {
		return nil, err
	}

	txRes, err := signedTx.Send()
	if err != nil {
		return nil, err
	}

	return txRes, nil
}
