package client

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rrossilli/glow/util"
)

// Retrieve cadence from file and replace imports with addresses from specified flow.json
func (c *GlowClient) CadenceFromFile(file string) (string, error) {
	p := path.Join(c.root, file)
	cdc, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}

	editedCdc := c.replaceImportAddresses(string(cdc))

	return editedCdc, nil
}

// Replace "0x" imports and file import paths in cadence with addresses from specified flow.json
func (c *GlowClient) replaceImportAddresses(cdc string) string {

	// Replace file import paths
	lines := strings.Split(string(cdc), "\n")
	for i, line := range lines {
		lines[i] = replaceFileImportPath(line)
	}
	newCdc := strings.Join(lines, "\n")

	// map to contract names
	contracts := c.FlowJSON.Contracts()
	var contractNames []string
	for cn := range contracts {
		contractNames = append(contractNames, cn)
	}

	// sort contract names by length (longest first) so that we replace the longest
	// contract names first. This is to avoid replacing a contract name that is a
	// substring of another contract name
	contractNamesSorted := util.SortStringsByCharacterLength(contractNames, false)

	// replace 0x imports
	for _, key := range contractNamesSorted {
		co := c.FlowJSON.Contract(key)
		newCdc = strings.Replace(
			newCdc,
			util.PrependHexPrefix(key),
			co.Address(c.network.Name),
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
