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

// TestTransferFlow verifies that transferring Flow tokens between accounts works as expected.
func TestTransferFlow(t *testing.T) {
	// Initialize a Glow client and retrieve the service account.
	c := client.NewGlowClient().Start()
	svc := c.SvcAcct

	// Generate a private key for the recipient account.
	privKey, err := c.NewPrivateKey(GENERATE_KEYS_SEED_PHRASE)
	require.NoError(t, err)
	assert.NotNil(t, privKey)

	// Create the recipient account.
	recipient, err := c.CreateAccount(privKey)
	require.NoError(t, err)
	assert.NotNil(t, recipient)

	// Define the transfer amount.
	amount, err := cadence.NewUFix64("10.0")
	require.NoError(t, err)

	// Execute a Flow token transfer from the service account to the recipient.
	txRes, err := c.NewTxFromFile(TxPath("flow_transfer"), svc).
		Args(amount, recipient.CadenceAddress()).
		SignAndSend()
	require.NoError(t, err)
	assert.NotNil(t, txRes)
	assert.NoError(t, txRes.Error)

	// Query the recipient's Flow token balance.
	result, err := c.NewScFromFile(ScPath("flow_balance"), recipient.CadenceAddress()).Exec()
	require.NoError(t, err)

	// Ensure the balance is as expected.
	balance, ok := result.ToGoValue().(uint64)
	require.True(t, ok, "balance should be a uint64")
	assert.Greater(t, balance, uint64(1), fmt.Sprintf("expected balance > 1, got %d", balance))
}
