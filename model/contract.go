package model

import "github.com/onflow/cadence"

// Contract struct as it typically appears in a flow.json
type Contract struct {
	Source  string            `json:"source"`
	Aliases map[string]string `json:"aliases"`
}

func (c Contract) Address(network string) string {
	return c.Aliases[network]
}

type ContractCdc struct {
	Contract Contract
	Name     string
	Cdc      string
}

// Contract cadence code as bytes
func (c ContractCdc) CdcBytes() []byte {
	return []byte(c.Cdc)
}

// Contract name as cadence string
func (c ContractCdc) NameAsCadenceString() cadence.String {
	return cadence.String(c.Name)
}
