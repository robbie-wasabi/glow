package actor

import (
	"fmt"

	"github.com/onflow/flow-go-sdk/crypto"

	"github.com/onflow/cadence"
	. "github.com/rrossilli/glow/client"
	. "github.com/rrossilli/glow/model"
	. "github.com/rrossilli/glow/util"
)

type Actor struct {
	Account Account
	Client  GlowClient
}

// Actor constructor
func NewActor(
	account Account,
	client GlowClient,
) Actor {
	return Actor{
		Account: account,
		Client:  client,
	}
}

// Create a new account on chain with a generic seed phrase and wrap in an Actor
func CreateDisposableActor(client GlowClient) (*Actor, error) {
	privKey, err := client.NewPrivateKey(DEFAULT_KEYS_SEED_PHRASE)
	if err != nil {
		return nil, err
	}

	acct, err := client.CreateAccount(privKey)
	if err != nil {
		return nil, err
	}

	return &Actor{
		Account: *acct,
		Client:  client,
	}, nil
}

// Create a new account on chain and wrap in an Actor
func CreateActor(
	privKey crypto.PrivateKey,
	client GlowClient,
) (*Actor, error) {
	acct, err := client.CreateAccount(privKey)
	if err != nil {
		return nil, err
	}

	return &Actor{
		Account: *acct,
		Client:  client,
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
func GetSvcActor(client GlowClient) Actor {
	account := client.FlowJSON.GetSvcAcct(client.GetNetwork())
	if IsEmpty(account) {
		panic(fmt.Sprintf("service account not found in flow.json: %s-svc", client.GetNetwork()))
	}
	return Actor{
		Account: account,
		Client:  client,
	}
}

// Get actor by name and current network
func GetActor(name string, client GlowClient) Actor {
	account := client.FlowJSON.GetAccount(fmt.Sprintf("%s-%s", client.GetNetwork(), name))
	if IsEmpty(account) {
		panic(fmt.Sprintf("account not found in flow.json: %s-svc", client.GetNetwork()))
	}
	return Actor{
		Account: account,
		Client:  client,
	}
}
