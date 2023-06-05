package client

import (
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-cli/pkg/flowkit"
	"github.com/onflow/flow-cli/pkg/flowkit/services"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"

	. "github.com/rrossilli/glow/model"
)

type Tx struct {
	script      flowkit.Script // TODO: flowkit's "script" type can be quite confusing since "script" is typically used to refer to reads from the chain
	payer       Account
	proposer    Account
	authorizers []Account
	client      *GlowClient // todo:
}

// Unsigned transaction contructor.
// Assumes that the proposer is also the gas payer and sole authorizer
func (c *GlowClient) NewTx(
	cdc []byte,
	proposer Account,
	args ...cadence.Value,
) *Tx {
	b := []byte(c.replaceImportAddresses(string(cdc)))
	return &Tx{
		script:   *flowkit.NewScript(b, args, ""),
		proposer: proposer,
		payer:    proposer,
		authorizers: []Account{
			proposer,
		},
		client: c,
	}
}

// Unsigned transaction contructor.
// Assumes that the proposer is also the gas payer and sole authorizer
func (c *GlowClient) NewTxFromString(
	cdc string,
	proposer Account,
	args ...cadence.Value,
) *Tx {
	b := []byte(c.replaceImportAddresses(cdc))
	return &Tx{
		script:   *flowkit.NewScript(b, args, ""),
		proposer: proposer,
		payer:    proposer,
		authorizers: []Account{
			proposer,
		},
		client: c,
	}
}

// Unsigned transaction contructor.
// Assumes that the proposer is also the gas payer and sole authorizer
func (c *GlowClient) NewTxFromFile(
	file string,
	proposer Account,
	args ...cadence.Value,
) *Tx {
	cdc, err := c.CadenceFromFile(file)
	if err != nil {
		panic(fmt.Sprintf("tx not found at: %s", file))
	}

	return &Tx{
		script:   *flowkit.NewScript([]byte(cdc), args, ""),
		proposer: proposer,
		payer:    proposer,
		authorizers: []Account{
			proposer,
		},
		client: c,
	}
}

// Specify args
func (t *Tx) Args(args ...cadence.Value) *Tx {
	t.script.Args = args
	return t
}

// Add arg to args
func (t *Tx) AddArg(arg cadence.Value) *Tx {
	t.script.Args = append(t.script.Args, arg)
	return t
}

// Specify payer
func (t *Tx) Payer(p Account) *Tx {
	t.payer = p
	t.authorizers = append(t.authorizers, p)
	return t
}

// Specify proposer
func (t *Tx) Proposer(p Account) *Tx {
	t.proposer = p
	t.authorizers = append(t.authorizers, p)
	return t
}

// Specify tx authorizers (typically unneeded)
func (t *Tx) Authorizers(a ...Account) *Tx {
	t.authorizers = a
	return t
}

// Append tx authorizer
func (t *Tx) AddAuthorizer(a Account) *Tx {
	t.authorizers = append(t.authorizers, a)
	return t
}

type SignedTx struct {
	flowTx *flowkit.Transaction
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

// Sign tx
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

	var txAddresses = services.NewTransactionAddresses(
		t.proposer.FlowAddress(),
		t.payer.FlowAddress(),
		FlowAddressesFromAccounts(t.authorizers),
	)

	// build flow tx
	flowTx, err := t.client.Services.Transactions.Build(
		txAddresses,
		0, // todo: which key?
		&t.script,
		t.client.gasLimit,
		t.client.network,
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
		flowTx: flowTx,
		client: t.client,
	}, err
}

// Send a signed Transaction
func (signedTx *SignedTx) Send() (*flow.TransactionResult, error) {
	_, res, err := signedTx.client.Services.Transactions.SendSigned(signedTx.flowTx)
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
