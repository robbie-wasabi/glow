package client

import (
	"crypto/rand"

	"github.com/onflow/flow-go-sdk/crypto"

	. "github.com/rrossilli/glow/util"
)

// Create new "crypto" private key from seed phrase.
func (c *GlowClient) NewPrivateKey(seedPhrase string) (crypto.PrivateKey, error) {
	seed := []byte(seedPhrase)
	_, err := rand.Read(seed)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.GeneratePrivateKey(c.SigAlgo, seed)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// Create new "crypto" private key from private key string.
func (c *GlowClient) NewPrivateKeyFromString(s string) (crypto.PrivateKey, error) {
	key, err := crypto.DecodePrivateKeyHex(c.SigAlgo, RemoveHexPrefix(s))
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Create new "crypto" public key from public key string.
func (c *GlowClient) NewPublicKeyFromString(s string) (crypto.PublicKey, error) {
	key, err := crypto.DecodePublicKeyHex(c.SigAlgo, RemoveHexPrefix(s))
	if err != nil {
		return nil, err
	}

	return key, nil
}
