package test

import (
	"testing"

	"github.com/onflow/cadence"
	. "github.com/rrossilli/glow/client"
	. "github.com/rrossilli/glow/util"

	. "github.com/smartystreets/goconvey/convey"
)

// Mint an NFT and transfer to collector account
// todo:
// test royalty
// test metadata views
func TestMintNFT(t *testing.T) {
	Convey("Create a client", t, func() {
		client := NewGlowClient().Start()
		minter := client.SvcAcct

		Convey("Set up royalty", func() {
			txRes, err := client.NewTxFromFile(
				TxPath("account_setup_royalty"),
				minter,
				cadence.Path{
					Domain:     "storage",
					Identifier: "flowTokenVault",
				},
			).SignAndSend()
			So(err, ShouldBeNil)
			So(txRes, ShouldNotBeNil)

			Convey("Mint NFTs", func() {
				txRes, err := client.NewTxFromFile(
					TxPath("nft_mint"),
					minter,
				).Args(
					minter.CadenceAddress(),
					cadence.String("name"),
					cadence.String("description"),
					cadence.String("thumbnail"),
					cadence.Array{
						Values: []cadence.Value{
							cadence.UFix64(100),
						},
					},
					cadence.Array{
						Values: []cadence.Value{
							cadence.String("royalty description"),
						},
					},
					cadence.Array{
						Values: []cadence.Value{
							minter.CadenceAddress(),
						},
					},
				).SignAndSend()
				So(err, ShouldBeNil)
				So(txRes, ShouldNotBeNil)

				Convey("Create a collector account", func() {
					collector, err := client.CreateDisposableAccount()
					So(err, ShouldBeNil)
					So(collector, ShouldNotBeNil)

					txRes, err := client.NewTxFromFile(
						TxPath("account_setup"),
						*collector,
					).SignAndSend()
					So(err, ShouldBeNil)
					So(txRes, ShouldNotBeNil)

					Convey("Transfer NFT from minter to collector", func() {
						txRes, err := client.NewTxFromFile(
							TxPath("nft_transfer"),
							minter,
							collector.CadenceAddress(),
							cadence.UInt64(0),
						).SignAndSend()
						So(err, ShouldBeNil)
						So(txRes, ShouldNotBeNil)

						nft, err := client.NewScFromFile(
							ScPath("nft_borrow"),
							collector.CadenceAddress(),
							cadence.UInt64(0),
						).Exec()
						So(err, ShouldBeNil)
						So(nft, ShouldNotBeNil)
					})
				})
			})
		})
	})
}
