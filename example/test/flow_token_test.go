package fcl

import (
	"fmt"

	"testing"

	. "github.com/rrossilli/glow/model"

	. "github.com/rrossilli/glow/client"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	GENERATE_KEYS_SEED_PHRASE = "elephant ears space cowboy octopus rodeo potato cannon pineapple"
)

// Get a specified Account's Flow token balance
func GetFlowTokenBalance(
	acctAddr cadence.Address,
	client GlowClient,
) (cadence.Value, error) {
	args := []cadence.Value{
		cadence.Address(acctAddr),
	}

	result, err := client.ExecScFromFile(scPath("flow_balance"), args)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Deposits Flow tokens into specified account
func DepositFlowTokens(
	amount cadence.UFix64,
	recipientAddr cadence.Address,
	proposer Account,
	client GlowClient,
) (*flow.TransactionResult, error) {
	args := []cadence.Value{
		amount, recipientAddr,
	}

	txRes, err := client.SignAndSendTxFromFile(
		scPath("flow_transfer"),
		args,
		proposer,
	)
	if err != nil {
		return nil, err
	}

	return txRes, nil
}

// Test Deposit Flow Tokens from service account into a newly created account
func TestDepositFlowTokens(t *testing.T) {
	Convey("Create a client", t, func() {

		// create and start new glow client
		client := NewGlowClient().Start()

		// get service account
		svcAcct := client.GetSvcAcct()

		Convey("Create a new account on the flow blockchain", func() {
			privKey, err := client.NewPrivateKey(GENERATE_KEYS_SEED_PHRASE)
			So(err, ShouldBeNil)
			So(privKey, ShouldNotBeNil)

			recipient, err := client.CreateAccount(
				privKey,
				svcAcct,
			)
			So(err, ShouldBeNil)
			So(recipient, ShouldNotBeNil)

			Convey("Deposit flow tokens into the account", func() {
				s := fmt.Sprintf("%v", "10.0")
				amount, err := cadence.NewUFix64(s)
				So(err, ShouldBeNil)

				res, err := DepositFlowTokens(
					amount,
					recipient.CadenceAddress(),
					svcAcct,
					*client,
				)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.Error, ShouldBeNil)

				Convey("Get flow token balance of account", func() {
					balance, err := GetFlowTokenBalance(recipient.CadenceAddress(), *client)
					So(err, ShouldBeNil)
					So(balance.ToGoValue(), ShouldBeGreaterThan, 1)
				})
			})
		})
	})
}
