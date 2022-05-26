package client

import (
	. "github.com/rrossilli/glow/model"
	. "github.com/rrossilli/glow/tmp"
	. "github.com/rrossilli/glow/util"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
)

// Get Account helper function
func (c *GlowClient) GetAccount(addr string) (*flow.Account, error) {
	a := flow.HexToAddress(addr)
	acct, err := c.Services.Accounts.Get(a)
	if err != nil {
		return nil, err
	}

	return acct, nil
}

// Create a new account on chain with a generic/unsafe seed phrase.
// These accounts are considered disposable as they have unsafe keys
func (c *GlowClient) CreateDisposableAccount() (*Account, error) {
	privKey, err := c.NewPrivateKey(DEFAULT_KEYS_SEED_PHRASE)
	if err != nil {
		return nil, err
	}

	acct, err := c.CreateAccount(privKey)
	if err != nil {
		return nil, err
	}

	return acct, err
}

// Create a new account on chain
func (c *GlowClient) CreateAccount(
	privKey crypto.PrivateKey,
) (*Account, error) {
	svcAcct := c.FlowJSON.GetSvcAcct(c.network)
	txRes, err := c.NewTx(
		[]byte(TX_CREATE_ACCOUNT),
		svcAcct,
		cadence.String(RemoveHexPrefix(privKey.PublicKey().String())),
	).SignAndSend()
	if err != nil {
		return nil, err
	}

	// fetch the address from the created account
	var address flow.Address
	if txRes.Status == flow.TransactionStatusSealed {
		for _, event := range txRes.Events {
			if event.Type == flow.EventAccountCreated {
				accountCreatedEvent := flow.AccountCreatedEvent(event)
				address = accountCreatedEvent.Address()
			}
		}
	}
	addrCdc := cadence.Address(address)

	a := NewAccount(
		addrCdc.String(),
		privKey.String(),
	)

	return &a, nil
}
