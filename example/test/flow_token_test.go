package test

import (
	"fmt"
	"testing"

	"github.com/onflow/cadence"
	"github.com/rrossilli/glow/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	GENERATE_KEYS_SEED_PHRASE = "elephant ears space cowboy octopus rodeo potato cannon pineapple"
)

func TestTransferFlow(t *testing.T) {
	client := client.NewGlowClient().Start()
	svc := client.SvcAcct

	// Create a new account on the flow blockchain
	privKey, err := client.NewPrivateKey(GENERATE_KEYS_SEED_PHRASE)
	require.Nil(t, err)
	assert.NotNil(t, privKey)

	recipient, err := client.CreateAccount(privKey)
	require.Nil(t, err)
	assert.NotNil(t, recipient)

	// Deposit flow tokens into the account
	s := fmt.Sprintf("%v", "10.0")
	amount, err := cadence.NewUFix64(s)
	require.Nil(t, err)

	txRes, err := client.NewTxFromFile(
		TxPath("flow_transfer"),
		svc,
	).Args(
		amount,
		recipient.CadenceAddress(),
	).SignAndSend()
	require.Nil(t, err)
	assert.NotNil(t, txRes)
	assert.Nil(t, txRes.Error)

	// Get flow token balance of account
	result, err := client.NewScFromFile(
		ScPath("flow_balance"),
		recipient.CadenceAddress(),
	).Exec()
	require.Nil(t, err)

	balance, ok := result.ToGoValue().(uint64)
	assert.True(t, ok)
	assert.Greater(t, balance, uint64(1))
}
