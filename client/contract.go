package client

import (
	"fmt"

	"github.com/rrossilli/glow/util"

	"github.com/rrossilli/glow/model"
)

// Get Contract by name
func (c *GlowClient) GetContractCdc(name string) model.ContractCdc {
	contract := c.FlowJSON.GetContract(name)
	if util.IsEmpty(contract) {
		panic(fmt.Sprintf("contract not found in flow.json: %s", name))
	}

	cdc, err := c.CadenceFromFile(contract.Source)
	if err != nil {
		panic(err)
	}

	return model.ContractCdc{
		Contract: contract,
		Name:     name,
		Cdc:      cdc,
	}
}
