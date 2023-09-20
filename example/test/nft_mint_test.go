package test

import (
	"testing"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/runtime/common"
	"github.com/rrossilli/glow/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMintNFT(t *testing.T) {
	client := client.NewGlowClient().Start()
	minter := client.SvcAcct

	// Set up royalty
	txRes, err := client.NewTxFromFile(
		TxPath("account_setup_royalty"),
		minter,
		cadence.Path{
			Domain:     common.PathDomainStorage,
			Identifier: "flowTokenVault",
		},
	).SignAndSend()
	require.Nil(t, err)
	assert.NotNil(t, txRes)

	// Mint NFTs
	txRes, err = client.NewTxFromFile(
		TxPath("nft_mint"),
		minter,
	).Args(
		minter.CadenceAddress(),
		cadence.String("name"),
		cadence.String("description"),
		cadence.String("thumbnail"),
		cadence.NewArray([]cadence.Value{cadence.UFix64(100)}),
		cadence.NewArray([]cadence.Value{cadence.String("royalty description")}),
		cadence.NewArray([]cadence.Value{minter.CadenceAddress()}),
	).SignAndSend()
	require.Nil(t, err)
	assert.NotNil(t, txRes)

	// Create a collector account
	collector, err := client.CreateDisposableAccount()
	require.Nil(t, err)
	assert.NotNil(t, collector)

	txRes, err = client.NewTxFromFile(
		TxPath("account_setup"),
		*collector,
	).SignAndSend()
	require.Nil(t, err)
	assert.NotNil(t, txRes)

	// Transfer NFT from minter to collector
	txRes, err = client.NewTxFromFile(
		TxPath("nft_transfer"),
		minter,
		collector.CadenceAddress(),
		cadence.UInt64(0),
	).SignAndSend()
	require.Nil(t, err)
	assert.NotNil(t, txRes)

	nft, err := client.NewScFromFile(
		ScPath("nft_borrow"),
		collector.CadenceAddress(),
		cadence.UInt64(0),
	).Exec()
	require.Nil(t, err)
	assert.NotNil(t, nft)
}
