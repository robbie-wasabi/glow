package test

import (
	"fmt"

	"testing"

	. "github.com/rrossilli/glow/client"
	. "github.com/rrossilli/glow/util"

	"github.com/onflow/cadence"
	. "github.com/smartystreets/goconvey/convey"
)

// Test Deposit Flow Tokens from service account into a newly created account
func TestTransferFlow(t *testing.T) {
	Convey("Create a client", t, func() {

		// create and start new glow client
		client := NewGlowClient().Start()

		// get service account
		svc := client.SvcAcct

		Convey("Create a new account on the flow blockchain", func() {
			privKey, err := client.NewPrivateKey(GENERATE_KEYS_SEED_PHRASE)
			So(err, ShouldBeNil)
			So(privKey, ShouldNotBeNil)

			recipient, err := client.CreateAccount(
				privKey,
			)
			So(err, ShouldBeNil)
			So(recipient, ShouldNotBeNil)

			Convey("Deposit flow tokens into the account", func() {
				s := fmt.Sprintf("%v", "10.0")
				amount, err := cadence.NewUFix64(s)
				So(err, ShouldBeNil)

				txRes, err := client.NewTxFromFile(
					TxPath("flow_transfer"),
					svc,
				).Args(
					amount,
					recipient.CadenceAddress(),
				).SignAndSend()
				So(err, ShouldBeNil)
				So(txRes, ShouldNotBeNil)
				So(txRes.Error, ShouldBeNil)

				Convey("Get flow token balance of account", func() {
					result, err := client.NewScFromFile(
						ScPath("flow_balance"),
						recipient.CadenceAddress(),
					).Exec()
					So(err, ShouldBeNil)
					So(result.ToGoValue().(uint64), ShouldBeGreaterThan, 1)
				})
			})
		})
	})
}
