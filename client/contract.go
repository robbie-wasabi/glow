package client

import (
	"fmt"
	"io/ioutil"
	"strings"

	. "github.com/rrossilli/glow/util"

	. "github.com/rrossilli/glow/model"
)

// Get Contract by name
func (c GlowClient) GetContract(name string) ContractCdc {
	contract := c.FlowJSON.GetContract(name)
	if IsEmpty(contract) {
		panic(fmt.Sprintf("contract not found in flow.json: %s", name))
	}

	// todo:
	source := RemoveFirstChar(contract.Source)

	cdc, err := c.GetContractCdc(source)
	if err != nil {
		panic(err)
	}

	return ContractCdc{
		Contract: contract,
		Name:     name,
		Cdc:      *cdc,
	}
}

// Retrieve cadence from file and replace imports with addresses from config
func (c GlowClient) GetContractCdc(file string) (*string, error) {
	p := c.root + file
	contractFile, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	cdc := c.replaceContractFilePaths(string(contractFile))

	return &cdc, nil
}

// Replace imported contract file paths
func (c GlowClient) replaceContractFilePaths(cdc string) string {
	keys := c.FlowJSON.ContractNamesSortedByLength(false)
	lines := strings.Split(string(cdc), "\n")
	for i, line := range lines {
		for _, key := range keys {
			co := c.FlowJSON.Contracts[key]
			if strings.Contains(line, fmt.Sprintf("import %v from", key)) {
				lines[i] = fmt.Sprintf("import %v from %v", key, co.Address(c.network))
			}
		}
	}
	return strings.Join(lines, "\n")
}
