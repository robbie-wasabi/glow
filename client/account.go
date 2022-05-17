package client

import (
	"fmt"
	"strings"

	. "github.com/rrossilli/glow/model"
	. "github.com/rrossilli/glow/util"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
)

// Get service account for current env
func (c GlowClient) GetSvcAcct() Account {
	account := c.FlowJSON.GetSvcAcct(c.env)
	if IsEmpty(account) {
		panic(fmt.Sprintf("account not found in flow.json: %s-svc", c.env))
	}
	return account
}

// Get Account by name and current env
func (c GlowClient) GetAccount(name string) Account {
	account := c.FlowJSON.GetAccount(fmt.Sprintf("%s-%s", c.env, name))
	if IsEmpty(account) {
		panic(fmt.Sprintf("account not found in flow.json: %s", name))
	}
	return account
}

// Get Accounts for current env
func (c GlowClient) Accounts() map[string]Account {
	accounts := map[string]Account{}
	for n, a := range c.FlowJSON.Accounts {
		if strings.Contains(n, c.env) {
			accounts[n] = a
		}
	}
	return accounts
}

// Get Account names for current env
func (c GlowClient) AccountNames() []string {
	keys := make([]string, 0, len(c.Accounts()))
	for k := range c.Accounts() {
		keys = append(keys, k)
	}
	return keys
}

// Sort accounts by predetermined emulator address order.
func (c GlowClient) AccountsSorted() []Account {
	var sorted []Account
	for _, o := range EMULATOR_ADDRESS_ORDER {
		for _, a := range c.Accounts() {
			if a.Address == o {
				sorted = append(sorted, a)
			}
		}
	}
	return sorted
}

// Sort account names by predetermined emulator address order.
func (c GlowClient) AccountNamesSorted() []string {
	var sorted []string
	for _, o := range EMULATOR_ADDRESS_ORDER {
		for n, a := range c.Accounts() {
			if a.Address == o {
				name := strings.ReplaceAll(n, fmt.Sprintf("%s-", c.env), "")
				sorted = append(sorted, name)
			}
		}
	}
	return sorted
}

// Create a new account on chain
func (c *GlowClient) CreateAccount(
	privKey crypto.PrivateKey,
	proposer Account,
) (*Account, error) {
	t := `
		transaction(publicKey: String) {
			prepare(signer: AuthAccount) {
				let account = AuthAccount(payer: signer)
				let accountKey = PublicKey(
					publicKey: publicKey.decodeHex(),
					signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
				)
				account.keys.add(
					publicKey: accountKey,
					hashAlgorithm: HashAlgorithm.SHA3_256,
					weight: 1000.0
				)
			}
		}
	`

	txRes, err := c.SignAndSendTx(
		t,
		[]cadence.Value{
			cadence.String(RemoveHexPrefix(privKey.PublicKey().String())),
		},
		proposer,
	)
	if err != nil {
		return nil, err
	}

	// fetch the address from the created account
	var address flow.Address
	if txRes.Status == flow.TransactionStatusSealed {
		for _, event := range txRes.Events {
			if event.Type == flow.EventAccountCreated {
				accountCreatedEvent := flow.AccountCreatedEvent(event)
				address = accountCreatedEvent.Address()
			}
		}
	}
	addrCdc := cadence.Address(address)

	a := NewAccount(
		addrCdc.String(),
		privKey.String(),
	)

	return &a, nil
}
