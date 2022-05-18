package model

import (
	"fmt"
	"sort"
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

// Names of accounts
func (f FlowJSON) AccountNames() []string {
	keys := make([]string, 0, len(f.Accounts))
	for k := range f.Accounts {
		keys = append(keys, k)
	}
	return keys
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
