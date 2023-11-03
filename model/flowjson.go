package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rrossilli/glow/consts"
	"github.com/rrossilli/glow/util"
)

// Maps to a standard flow.json
type FlowJSON struct {
	data struct {
		Emulator    interface{}           `json:"emulators"` // todo
		Contracts   map[string]Contract   `json:"contracts"`
		Networks    map[string]string     `json:"networks"`
		Accounts    map[string]Account    `json:"accounts"`
		Deployments map[string]Deployment `json:"deployments"`
	}
}

// construct FlowJSON from json bytes
// func NewFlowJSON(b []byte) (FlowJSON, error) {
// 	var f FlowJSON
// 	err := util.UnmarshalJSON(b, &f)
// 	if err != nil {
// 		return f, err
// 	}
// 	return f, nil
// }

func (f FlowJSON) FromBytes(b []byte) (FlowJSON, error) {
	err := json.Unmarshal(b, &f)
	if err != nil {
		return f, err
	}
	return f, nil
}

func (f FlowJSON) Contract(name string) Contract {
	return f.data.Contracts[name]
}

// Get contracts
func (f FlowJSON) Contracts() map[string]Contract {
	return f.data.Contracts
}

// Names of contracts
// func (f FlowJSON) ContractNames() []string {
// 	keys := make([]string, 0, len(f.data.Contracts))
// 	for k := range f.data.Contracts {
// 		keys = append(keys, k)
// 	}
// 	return keys
// }

// Sort contract names by length. Helpful for replacing import Addresses
// in scripts
// func (f FlowJSON) ContractNamesSortedByLength(asc bool) []string {
// 	keys := make([]string, 0, len(f.data.Contracts))
// 	for k := range f.data.Contracts {
// 		keys = append(keys, k)
// 	}
// 	sort.SliceStable(keys, func(i, j int) bool {
// 		if asc {
// 			return len(keys[i]) < len(keys[j])
// 		} else {
// 			return len(keys[i]) > len(keys[j])
// 		}
// 	})
// 	return keys
// }

// Get service account (account with "svc" suffix)
func (f FlowJSON) ServiceAccount(network string) Account {
	return f.Account(fmt.Sprintf("%s-svc", network))
}

// Get account by name
func (f FlowJSON) Account(name string) Account {
	account := f.data.Accounts[name]
	return account
}

// Get Accounts for network
func (f FlowJSON) Accounts(network string) map[string]Account {
	accounts := map[string]Account{}
	for n, a := range f.data.Accounts {
		if strings.Contains(n, network) {
			accounts[n] = a
		}
	}
	return accounts
}

// Sort accounts for network by predetermined emulator address order.
func (f FlowJSON) AccountsSorted() []Account {
	var sorted []Account
	for _, o := range consts.EMULATOR_ADDRESS_ORDER {
		for _, a := range f.data.Accounts {
			if util.PrependHexPrefix(a.Address) == util.PrependHexPrefix(o) {
				sorted = append(sorted, a)
			}
		}
	}
	return sorted
}

// Helper function to sort account names for network by predetermined emulator address order.
func (f FlowJSON) AccountNamesSorted(network string) []string {
	var sorted []string
	for _, o := range consts.EMULATOR_ADDRESS_ORDER {
		fmt.Printf("f.data.Accounts: %v\n", f.data.Accounts)
		for n, a := range f.data.Accounts {
			if util.PrependHexPrefix(a.Address) == util.PrependHexPrefix(o) {
				// name := strings.ReplaceAll(n, fmt.Sprintf("%s-", network), "")
				// sorted = append(sorted, name)
				sorted = append(sorted, n)
			}
		}
	}
	return sorted
}

func (f FlowJSON) Deployment(network string) Deployment {
	deployment := f.data.Deployments[network]
	return deployment
}
