package client

import (
	"fmt"

	. "github.com/rrossilli/glow/util"

	. "github.com/rrossilli/glow/model"
)

// Get Contract by name
func (c *GlowClient) GetContractCdc(name string) ContractCdc {
	contract := c.FlowJSON.GetContract(name)
	if IsEmpty(contract) {
		panic(fmt.Sprintf("contract not found in flow.json: %s", name))
	}

	cdc, err := c.CadenceFromFile(contract.Source)
	if err != nil {
		panic(err)
	}

	return ContractCdc{
		Contract: contract,
		Name:     name,
		Cdc:      cdc,
	}
}
