package actor

import (
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk/crypto"

	. "github.com/rrossilli/glow/client"
	. "github.com/rrossilli/glow/util"
)

type WrappedGlowClient struct {
	Client *GlowClient
}

func NewWrappedGlowClient(client *GlowClient) WrappedGlowClient {
	return WrappedGlowClient{
		Client: client,
	}
}

// Create a new account on chain with a generic seed phrase and wrap in an Actor
func (w WrappedGlowClient) CreateDisposableActor() (*Actor, error) {
	privKey, err := w.Client.NewPrivateKey(DEFAULT_KEYS_SEED_PHRASE)
	if err != nil {
		return nil, err
	}

	acct, err := w.Client.CreateAccount(privKey)
	if err != nil {
		return nil, err
	}

	return &Actor{
		Account: *acct,
		Client:  w.Client,
	}, nil
}

// Create a new account on chain and wrap in an Actor
func (w WrappedGlowClient) CreateActor(
	privKey crypto.PrivateKey,
) (*Actor, error) {
	acct, err := w.Client.CreateAccount(privKey)
	if err != nil {
		return nil, err
	}

	return &Actor{
		Account: *acct,
		Client:  w.Client,
	}, nil
}

// Get service account actor for current network
func (w WrappedGlowClient) GetSvcActor() Actor {
	account := w.Client.FlowJSON.GetSvcAcct(w.Client.GetNetwork())
	if IsEmpty(account) {
		panic(fmt.Sprintf("service account not found in flow.json: %s-svc", w.Client.GetNetwork()))
	}
	return Actor{
		Account: account,
		Client:  w.Client,
	}
}

// Get actor by name and current network
func (w WrappedGlowClient) GetActor(name string) Actor {
	account := w.Client.FlowJSON.GetAccount(fmt.Sprintf("%s-%s", w.Client.GetNetwork(), name))
	if IsEmpty(account) {
		panic(fmt.Sprintf("account not found in flow.json: %s-svc", w.Client.GetNetwork()))
	}
	return Actor{
		Account: account,
		Client:  w.Client,
	}
}

// Execute a Script from a string
func (w WrappedGlowClient) ExecSc(
	cdc string,
	args ...cadence.Value,
) (cadence.Value, error) {
	result, err := w.Client.ExecSc(cdc, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Execute a Script at a specified file path
func (w WrappedGlowClient) ExecScFromFile(
	file string,
	args ...cadence.Value,
) (cadence.Value, error) {
	sc, err := w.Client.NewScFromFile(file)
	if err != nil {
		return nil, err
	}

	res, err := w.Client.ExecSc(sc.Script, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
