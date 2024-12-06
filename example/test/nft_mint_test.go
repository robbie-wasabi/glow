package test

import (
	"testing"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/runtime/common"
	"github.com/rrossilli/glow/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMintNFT verifies NFT minting and transferring between accounts.
func TestMintNFT(t *testing.T) {
	c := client.NewGlowClient().Start()
	minter := c.SvcAcct

	// Set up a royalty vault in the minter’s account.
	txRes, err := c.NewTxFromFile(
		TxPath("account_setup_royalty"),
		minter,
		cadence.Path{
			Domain:     common.PathDomainStorage,
			Identifier: "flowTokenVault",
		},
	).SignAndSend()
	require.NoError(t, err)
	assert.NotNil(t, txRes)

	// Mint an NFT.
	txRes, err = c.NewTxFromFile(
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
	require.NoError(t, err)
	assert.NotNil(t, txRes)

	// Create a disposable collector account and set it up for NFTs.
	collector, err := c.CreateDisposableAccount()
	require.NoError(t, err)
	assert.NotNil(t, collector)

	txRes, err = c.NewTxFromFile(
		TxPath("account_setup"),
		*collector,
	).SignAndSend()
	require.NoError(t, err)
	assert.NotNil(t, txRes)

	// Transfer the newly minted NFT to the collector.
	txRes, err = c.NewTxFromFile(
		TxPath("nft_transfer"),
		minter,
		collector.CadenceAddress(),
		cadence.UInt64(0),
	).SignAndSend()
	require.NoError(t, err)
	assert.NotNil(t, txRes)

	// Verify that the collector can borrow the NFT.
	nft, err := c.NewScFromFile(
		ScPath("nft_borrow"),
		collector.CadenceAddress(),
		cadence.UInt64(0),
	).Exec()
	require.NoError(t, err)
	assert.NotNil(t, nft)
}
