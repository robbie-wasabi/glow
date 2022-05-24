package client

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	. "github.com/rrossilli/glow/util"
)

// Retrieve cadence from file and replace imports with addresses from specified flow.json
func (c GlowClient) CadenceFromFile(file string) (string, error) {
	p := path.Join(c.root, file)
	cdc, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}

	editedCdc := c.replaceImportAddresses(string(cdc))

	return editedCdc, nil
}

// Replaces "0x" imports and import paths in cadence with addresses from specified flow.json
func (c GlowClient) replaceImportAddresses(cdc string) string {
	keys := c.FlowJSON.ContractNamesSortedByLength(false)

	// replace 0x imports i.e. "import Contract from 0xContract"
	for _, key := range keys {
		co := c.FlowJSON.Contracts[key]
		cdc = strings.Replace(
			cdc,
			PrependHexPrefix(key),
			co.Address(c.network),
			-1,
		)
	}

	// replace import paths i.e. "import Contract from '../contracts/Contract.cdc'"
	lines := strings.Split(string(cdc), "\n")
	for i, line := range lines {
		for _, key := range keys {
			co := c.FlowJSON.Contracts[key]
			if strings.Contains(line, fmt.Sprintf("import %v from", key)) {
				lines[i] = fmt.Sprintf("import %v from %v", key, co.Address(c.network))
			}
		}
	}

	return cdc
}
