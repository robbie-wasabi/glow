package test

import (
	"encoding/hex"
	"testing"

	"github.com/onflow/cadence"
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
			// sign data
			data := "test"
			signedData, err := svc.SignMessage([]byte(data), crypto.SHA3_256)

			// verify signed data on chain
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
