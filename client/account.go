package client

import (
	"fmt"
	"strings"

	. "github.com/rrossilli/glow/model"
	. "github.com/rrossilli/glow/tmp"
	. "github.com/rrossilli/glow/util"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
)

// Get service account for current network
func (c GlowClient) GetSvcAcct() Account {
	account := c.FlowJSON.GetSvcAcct(c.network)
	if IsEmpty(account) {
		panic(fmt.Sprintf("account not found in flow.json: %s-svc", c.network))
	}
	return account
}

// Get Account by name and current network
func (c GlowClient) GetAccount(name string) Account {
	account := c.FlowJSON.GetAccount(fmt.Sprintf("%s-%s", c.network, name))
	if IsEmpty(account) {
		panic(fmt.Sprintf("account not found in flow.json: %s", name))
	}
	return account
}

// Get Accounts for current network
func (c GlowClient) Accounts() map[string]Account {
	accounts := map[string]Account{}
	for n, a := range c.FlowJSON.Accounts {
		if strings.Contains(n, c.network) {
			accounts[n] = a
		}
	}
	return accounts
}

// Get Account names for current network
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
				name := strings.ReplaceAll(n, fmt.Sprintf("%s-", c.network), "")
				sorted = append(sorted, name)
			}
		}
	}
	return sorted
}

// Create a new account on chain with a generic/unsafe seed phrase.
// These accounts are considered disposable as they have unsafe keys
func (c GlowClient) CreateDisposableAccount() (*Account, error) {
	privKey, err := c.NewPrivateKey(DEFAULT_KEYS_SEED_PHRASE)
	if err != nil {
		return nil, err
	}

	acct, err := c.CreateAccount(privKey)
	if err != nil {
		return nil, err
	}

	return acct, err
}

// Create a new account on chain
func (c GlowClient) CreateAccount(
	privKey crypto.PrivateKey,
) (*Account, error) {
	svcAcct := c.GetSvcAcct()
	txRes, err := c.NewTx(
		[]byte(TX_CREATE_ACCOUNT),
		svcAcct,
		cadence.String(RemoveHexPrefix(privKey.PublicKey().String())),
	).SignAndSend()
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
