package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rrossilli/glow/consts"
	"github.com/rrossilli/glow/util"
)

// FlowJSON maps a standard flow.json structure.
type FlowJSON struct {
	data struct {
		Emulator    interface{}           `json:"emulators"`
		Contracts   map[string]Contract   `json:"contracts"`
		Networks    map[string]string     `json:"networks"`
		Accounts    map[string]Account    `json:"accounts"`
		Deployments map[string]Deployment `json:"deployments"`
	}
}

// FromBytes unmarshals FlowJSON from JSON bytes.
func (f FlowJSON) FromBytes(b []byte) (FlowJSON, error) {
	err := json.Unmarshal(b, &f)
	return f, err
}

// Contract returns the named contract.
func (f FlowJSON) Contract(name string) Contract {
	return f.data.Contracts[name]
}

// Contracts returns all defined contracts.
func (f FlowJSON) Contracts() map[string]Contract {
	return f.data.Contracts
}

// ServiceAccount returns the service account for the given network.
func (f FlowJSON) ServiceAccount(network string) Account {
	return f.Account(fmt.Sprintf("%s-svc", network))
}

// Account returns the account by name.
func (f FlowJSON) Account(name string) Account {
	return f.data.Accounts[name]
}

// Accounts returns all accounts for a given network.
func (f FlowJSON) Accounts(network string) map[string]Account {
	accounts := map[string]Account{}
	for n, a := range f.data.Accounts {
		if strings.Contains(n, network) {
			accounts[n] = a
		}
	}
	return accounts
}

// AccountsSorted returns emulator accounts in the predefined order.
func (f FlowJSON) AccountsSorted() []Account {
	var sorted []Account
	for _, addr := range consts.EMULATOR_ADDRESS_ORDER {
		for _, a := range f.data.Accounts {
			if util.PrependHexPrefix(a.Address) == util.PrependHexPrefix(addr) {
				sorted = append(sorted, a)
			}
		}
	}
	return sorted
}

// AccountNamesSorted returns emulator account names in the predefined order.
func (f FlowJSON) AccountNamesSorted(network string) []string {
	var sorted []string
	for _, addr := range consts.EMULATOR_ADDRESS_ORDER {
		fmt.Printf("f.data.Accounts: %v\n", f.data.Accounts)
		for n, a := range f.data.Accounts {
			if util.PrependHexPrefix(a.Address) == util.PrependHexPrefix(addr) {
				sorted = append(sorted, n)
			}
		}
	}
	return sorted
}

// Deployment returns the deployment configuration for the given network.
func (f FlowJSON) Deployment(network string) Deployment {
	return f.data.Deployments[network]
}
