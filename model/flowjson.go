package model

import (
	"fmt"
	"sort"
	"strings"

	. "github.com/rrossilli/glow/consts"
)

// Maps to a standard flow.json
type FlowJSON struct {
	Emulator    interface{}                    `json:"emulators"` // todo
	Contracts   map[string]Contract            `json:"contracts"`
	Networks    map[string]string              `json:"networks"`
	Accounts    map[string]Account             `json:"accounts"`
	Deployments map[string]map[string][]string `json:"deployments"`
}

// Get contract by name
func (f FlowJSON) GetContract(name string) Contract {
	return f.Contracts[name]
}

// Names of contracts
func (f FlowJSON) ContractNames() []string {
	keys := make([]string, 0, len(f.Contracts))
	for k := range f.Contracts {
		keys = append(keys, k)
	}
	return keys
}

// Sort contract names by length. Helpful for replacing import Addresses
// in scripts
func (f FlowJSON) ContractNamesSortedByLength(asc bool) []string {
	keys := make([]string, 0, len(f.Contracts))
	for k := range f.Contracts {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		if asc {
			return len(keys[i]) < len(keys[j])
		} else {
			return len(keys[i]) > len(keys[j])
		}
	})
	return keys
}

// Get account with "svc" suffix
func (f FlowJSON) GetSvcAcct(network string) Account {
	return f.GetAccount(fmt.Sprintf("%s-svc", network))
}

// Get account by name
func (f FlowJSON) GetAccount(name string) Account {
	account := f.Accounts[name]
	return account
}

// Get Accounts for network
func (f FlowJSON) GetAccounts(network string) map[string]Account {
	accounts := map[string]Account{}
	for n, a := range f.Accounts {
		if strings.Contains(n, network) {
			accounts[n] = a
		}
	}
	return accounts
}

// Names of accounts for network
func (f FlowJSON) AccountNames(network string) []string {
	accounts := f.GetAccounts(network)
	keys := make([]string, 0, len(accounts))
	for k := range f.Accounts {
		keys = append(keys, k)
	}
	return keys
}

// Sort accounts for network by predetermined emulator address order.
func (f FlowJSON) AccountsSorted() []Account {
	var sorted []Account
	for _, o := range EMULATOR_ADDRESS_ORDER {
		for _, a := range f.Accounts {
			if a.Address == o {
				sorted = append(sorted, a)
			}
		}
	}
	return sorted
}

// Sort account names for network by predetermined emulator address order.
func (f FlowJSON) AccountNamesSorted(network string) []string {
	var sorted []string
	for _, o := range EMULATOR_ADDRESS_ORDER {
		for n, a := range f.Accounts {
			if a.Address == o {
				// name := strings.ReplaceAll(n, fmt.Sprintf("%s-", network), "")
				// sorted = append(sorted, name)
				sorted = append(sorted, n)
			}
		}
	}
	return sorted
}

// Get deployment by name
func (f FlowJSON) GetDeployment(network string) map[string][]string {
	deployment := f.Deployments[network]
	return deployment
}

// Get deployment contracts by account name
func (f FlowJSON) GetAccountDeployment(network, name string) []string {
	deployment := f.Deployments[network]
	return deployment[fmt.Sprintf("%s-%s", network, name)]
}
