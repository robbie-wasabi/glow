package client

import (
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-cli/pkg/flowkit"
	"github.com/onflow/flow-cli/pkg/flowkit/util"
	"github.com/onflow/flow-go-sdk"
)

type Sc struct {
	script *flowkit.Script
	client *GlowClient
}

// Create new script
func (c *GlowClient) NewSc(bytes []byte, args ...cadence.Value) *Sc {
	b := []byte(c.replaceImportAddresses(string(bytes)))
	return &Sc{
		script: flowkit.NewScript(b, args, ""),
		client: c,
	}
}

// Create new script from string
func (c *GlowClient) NewScFromString(cdc string, args ...cadence.Value) *Sc {
	b := []byte(c.replaceImportAddresses(cdc))
	return &Sc{
		script: flowkit.NewScript(b, args, ""),
		client: c,
	}
}

// Create new script from file
func (c *GlowClient) NewScFromFile(file string, args ...cadence.Value) *Sc {
	cdc, err := c.CadenceFromFile(file)
	if err != nil {
		panic(fmt.Sprintf("sc not found at: %s", file))
	}
	b := []byte(c.replaceImportAddresses(cdc))
	return &Sc{
		script: flowkit.NewScript(b, args, ""),
		client: c,
	}
}

// Specify args
func (sc *Sc) Args(args ...cadence.Value) *Sc {
	sc.script.Args = args
	return sc
}

// Add arg to args
func (sc *Sc) AddArg(arg cadence.Value) *Sc {
	sc.script.Args = append(sc.script.Args, arg)
	return sc
}

// Execute script
func (sc *Sc) Exec() (cadence.Value, error) {
	// necessary or throws null pointer exception
	query := util.ScriptQuery{
		ID:     flow.EmptyID,
		Height: 0,
	}
	result, err := sc.client.Services.Scripts.Execute(sc.script, sc.client.network, &query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Execute script with query
func (sc *Sc) ExecWithQuery(query *util.ScriptQuery) (cadence.Value, error) {
	result, err := sc.client.Services.Scripts.Execute(sc.script, sc.client.network, query)
	if err != nil {
		return nil, err
	}

	return result, nil
}
