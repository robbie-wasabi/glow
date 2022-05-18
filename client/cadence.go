package client

import (
	"io/ioutil"
	"strings"

	. "github.com/rrossilli/glow/util"
)

// Retrieve cadence from file and replace imports with addresses from config
func (c GlowClient) CadenceFromFile(file string) (string, error) {

	p := c.root + file
	txFile, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}

	cdc := c.replaceImportAddresses(string(txFile))

	return cdc, nil
}

// Similar to FCL Config, replaces imports in cadence with addresses from config
func (c GlowClient) replaceImportAddresses(cdc string) string {
	keys := c.FlowJSON.ContractNamesSortedByLength(false)
	for _, key := range keys {
		co := c.FlowJSON.Contracts[key]
		cdc = strings.Replace(
			cdc,
			PrependHexPrefix(key),
			co.Address(c.network),
			-1,
		)
	}

	return cdc
}
