package model

import "github.com/onflow/cadence"

// Contract represents a contract entry in flow.json.
type Contract struct {
	Source  string            `json:"source"`
	Aliases map[string]string `json:"aliases"`
}

// Address returns the contract's deployment address on the given network.
func (c Contract) Address(network string) string {
	return c.Aliases[network]
}

// ContractCdc wraps a Contract with its name and Cadence code.
type ContractCdc struct {
	Contract Contract
	Name     string
	Cdc      string
}

// CdcBytes returns the raw Cadence code as a byte slice.
func (c ContractCdc) CdcBytes() []byte {
	return []byte(c.Cdc)
}

// NameAsCadenceString returns the contract name as a Cadence string.
func (c ContractCdc) NameAsCadenceString() cadence.String {
	return cadence.String(c.Name)
}
