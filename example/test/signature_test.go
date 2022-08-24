package test

import (
	"encoding/hex"
	"testing"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	. "github.com/rrossilli/glow/client"
	. "github.com/rrossilli/glow/util"

	. "github.com/smartystreets/goconvey/convey"
)

// Test signing arbitrary data
func TestSignData(t *testing.T) {
	Convey("Create a client", t, func() {

		// create and start new glow client
		client := NewGlowClient().Start()

		// get service account
		svc := client.SvcAcct

		Convey("Sign data", func() {
			// create crypto signer
			signer, err := crypto.NewInMemorySigner(svc.CryptoPrivateKey(), crypto.SHA3_256)
			So(err, ShouldBeNil)

			// sign data
			data := "test"
			signedData, err := flow.SignUserMessage(signer, []byte(data))

			// verify signed data in cadence
			res, err := client.NewScFromFile(
				ScPath("sig_verify"),
				cadence.String(RemoveHexPrefix(svc.CryptoPublicKey().String())),
				cadence.String(hex.EncodeToString(signedData)),
				cadence.String(data),
			).Exec()
			So(err, ShouldBeNil)
			So(res.ToGoValue().(bool), ShouldBeTrue)
		})
	})
}
