package actor

import (
	"github.com/onflow/cadence"
	. "github.com/rrossilli/glow/client"
	. "github.com/rrossilli/glow/model"
)

type Actor struct {
	Account Account
	Client  *GlowClient
}

// Actor constructor
func NewActor(
	account Account,
	client *GlowClient,
) Actor {
	return Actor{
		Account: account,
		Client:  client,
	}
}

// Unsigned transaction contstructor.
// Assumes that the actor is also the gas payer and sole authorizer
func (a Actor) NewTx(
	cdc []byte,
	args ...cadence.Value,
) Tx {
	tx := a.Client.NewTx(
		cdc,
		args,
		a.Account,
	)

	return tx
}

// Unsigned transaction contstructor.
// Assumes that the actor is also the gas payer and sole authorizer
func (a Actor) NewTxFromString(
	cdc string,
	args ...cadence.Value,
) Tx {
	tx := a.Client.NewTx(
		[]byte(cdc),
		args,
		a.Account,
	)

	return tx
}

// Unsigned transaction contstructor.
// Assumes that the actor is also the gas payer and sole authorizer
func (a Actor) NewTxFromFile(
	file string,
	args ...cadence.Value,
) Tx {
	tx := a.Client.NewTxFromFile(
		file,
		args,
		a.Account,
	)

	return tx
}
