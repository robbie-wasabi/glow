package client

import (
	"fmt"

	"github.com/onflow/flow-go-sdk/crypto"

	"github.com/onflow/cadence"
	. "github.com/rrossilli/glow/model"
	. "github.com/rrossilli/glow/util"
)

type Actor struct {
	Account Account
	Client  GlowClient
}

// Actor constructor
func (c GlowClient) NewActor(
	account Account,
) Actor {
	return Actor{
		Account: account,
		Client:  c,
	}
}

// Create a new account on chain with a generic seed phrase and wrap in an Actor
func (c GlowClient) CreateDisposableActor() (*Actor, error) {
	privKey, err := c.NewPrivateKey(DEFAULT_KEYS_SEED_PHRASE)
	if err != nil {
		return nil, err
	}

	acct, err := c.CreateAccount(privKey)
	if err != nil {
		return nil, err
	}

	return &Actor{
		Account: *acct,
		Client:  c,
	}, nil
}

// Create a new account on chain and wrap in an Actor
func (c GlowClient) CreateActor(
	privKey crypto.PrivateKey,
) (*Actor, error) {
	acct, err := c.CreateAccount(privKey)
	if err != nil {
		return nil, err
	}

	return &Actor{
		Account: *acct,
		Client:  c,
	}, nil
}

// Unsigned transaction contstructor.
// Assumes that the actor is also the gas payer and sole authorizer
func (a Actor) NewTx(
	cdc []byte,
	args ...cadence.Value,
) Tx {
	tx := NewTx(
		cdc,
		args,
		a.Account,
		a.Client,
	)

	return tx
}

// Unsigned transaction contstructor.
// Assumes that the actor is also the gas payer and sole authorizer
func (a Actor) NewTxFromString(
	cdc string,
	args ...cadence.Value,
) Tx {
	tx := NewTx(
		[]byte(cdc),
		args,
		a.Account,
		a.Client,
	)

	return tx
}

// Unsigned transaction contstructor.
// Assumes that the actor is also the gas payer and sole authorizer
func (a Actor) NewTxFromFile(
	file string,
	args ...cadence.Value,
) Tx {
	tx := NewTxFromFile(
		file,
		args,
		a.Account,
		a.Client,
	)

	return tx
}

// Get service account actor for current network
func (c GlowClient) GetSvcActor() Actor {
	account := c.FlowJSON.GetSvcAcct(c.network)
	if IsEmpty(account) {
		panic(fmt.Sprintf("service account not found in flow.json: %s-svc", c.network))
	}
	return Actor{
		Account: account,
		Client:  c,
	}
}

// Get actor by name and current network
func (c GlowClient) GetActor(name string) Actor {
	account := c.FlowJSON.GetAccount(fmt.Sprintf("%s-%s", c.network, name))
	if IsEmpty(account) {
		panic(fmt.Sprintf("account not found in flow.json: %s-svc", c.network))
	}
	return Actor{
		Account: account,
		Client:  c,
	}
}
