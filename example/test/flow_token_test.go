package fcl

import (
	"fmt"

	"testing"

	. "github.com/rrossilli/glow/actor"
	. "github.com/rrossilli/glow/client"
	. "github.com/rrossilli/glow/util"

	"github.com/onflow/cadence"
	. "github.com/smartystreets/goconvey/convey"
)

// Test Deposit Flow Tokens from service account into a newly created account
func TestDepositFlowTokens(t *testing.T) {
	Convey("Create a client", t, func() {

		// create and start new glow client
		client := NewGlowClient().Start()

		// get service account
		svcAcct := GetSvcActor(*client)

		Convey("Create a new account on the flow blockchain", func() {
			recipient, err := CreateDisposableActor(*client)
			So(err, ShouldBeNil)
			So(recipient, ShouldNotBeNil)

			Convey("Deposit flow tokens into the account", func() {
				s := fmt.Sprintf("%v", "10.0")
				amount, err := cadence.NewUFix64(s)
				So(err, ShouldBeNil)

				txRes, err := svcAcct.
					NewTxFromFile(TxPath("flow_transfer")).
					Args(
						amount,
						recipient.Account.CadenceAddress(),
					).
					SignAndSend()
				So(err, ShouldBeNil)
				So(txRes, ShouldNotBeNil)
				So(txRes.Error, ShouldBeNil)

				Convey("Get flow token balance of account", func() {
					result, err := client.ExecScFromFile(
						ScPath("flow_balance"),
						recipient.Account.CadenceAddress(),
					)
					So(err, ShouldBeNil)
					So(result.ToGoValue(), ShouldBeGreaterThan, 1)
				})
			})
		})
	})
}
