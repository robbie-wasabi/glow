package client

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	. "github.com/rrossilli/glow/util"
)

// Retrieve cadence from file and replace imports with addresses from specified flow.json
func (c *GlowClient) CadenceFromFile(file string) (string, error) {
	p := path.Join(c.root, file)
	cdc, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}

	editedCdc := c.replaceImportAddresses(string(cdc))

	return editedCdc, nil
}

// Replaces "0x" imports and file import paths in cadence with addresses from specified flow.json
func (c *GlowClient) replaceImportAddresses(cdc string) string {
	lines := strings.Split(string(cdc), "\n")
	for i, line := range lines {
		lines[i] = replaceFileImportPath(line)
	}
	newCdc := strings.Join(lines, "\n")

	// replace 0x imports i.e. "import Contract from 0xContract"
	keys := c.FlowJSON.ContractNamesSortedByLength(false)
	for _, key := range keys {
		co := c.FlowJSON.Contracts[key]
		newCdc = strings.Replace(
			newCdc,
			PrependHexPrefix(key),
			co.Address(c.network),
			-1,
		)
	}

	return newCdc
}

// Replace relative file import path with "0x path"
// i.e. import NonFungibleToken from "./NonFungibleToken.cdc" as ...from 0xNonFungibleToken
func replaceFileImportPath(s string) string {
	if strings.Contains(s, "import") &&
		strings.Contains(s, "from") &&
		strings.Contains(s, `.cdc`) {
		fields := strings.Fields(s)
		fields[3] = fmt.Sprintf("0x%s", fields[1])
		return strings.Join(fields, " ")
	}
	return s
}
