package client

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"

	. "github.com/rrossilli/glow/util"

	. "github.com/rrossilli/glow/model"
)

// Get Contract by name
func (c GlowClient) GetContract(name string) Contract {
	contract := c.FlowJSON.GetContract(name)
	if IsEmpty(contract) {
		panic(fmt.Sprintf("contract not found in flow.json: %s", name))
	}
	return contract
}

// Retrieve cadence from file and replace imports with addresses from config
func (c GlowClient) Contract(filePath string) (*string, error) {
	path := ROOT + filePath // we use absolute paths
	contractFile, err := ioutil.ReadFile(path)
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
				lines[i] = fmt.Sprintf("import %v from %v", key, co.Address(c.env))
			}
		}
	}
	return strings.Join(lines, "\n")
}

// Deploy contract to account
func (c *GlowClient) DeployContract(
	contractName cadence.String,
	contractFilePath string,
	proposer Account,
) (*flow.TransactionResult, error) {
	contract, err := c.Contract(contractFilePath)
	if err != nil {
		return nil, err
	}

	t := `
		transaction(name: String, code: String) {
			prepare(signer: AuthAccount) {
				signer.contracts.add(name: name, code: code.decodeHex())
			}
		}
	`

	txRes, err := c.SignAndSendTx(
		t,
		proposer,
		contractName,
		cadence.String(hex.EncodeToString([]byte(*contract))),
	)
	if err != nil {
		return nil, err
	}

	return txRes, nil
}
