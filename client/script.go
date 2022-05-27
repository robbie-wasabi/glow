package client

import (
	"fmt"

	"github.com/onflow/cadence"
)

type Sc struct {
	cdc    []byte
	args   []cadence.Value
	client *GlowClient
}

// Create new script
func (c *GlowClient) NewSc(bytes []byte, args ...cadence.Value) *Sc {
	b := []byte(c.replaceImportAddresses(string(bytes)))
	return &Sc{
		cdc:    b,
		args:   args,
		client: c,
	}
}

// Create new script from string
func (c *GlowClient) NewScFromString(cdc string, args ...cadence.Value) *Sc {
	b := []byte(c.replaceImportAddresses(cdc))
	return &Sc{
		cdc:    b,
		args:   args,
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
		cdc:    b,
		args:   args,
		client: c,
	}
}

// Specify args
func (sc *Sc) Args(args ...cadence.Value) *Sc {
	sc.args = args
	return sc
}

// Add arg to args
func (sc *Sc) AddArg(arg cadence.Value) *Sc {
	sc.args = append(sc.args, arg)
	return sc
}

// Execute script
func (sc *Sc) Exec() (cadence.Value, error) {
	// we don't need to pass the file name as we have a different strategy to replace imports
	result, err := sc.client.Services.Scripts.Execute(sc.cdc, sc.args, "", sc.client.network)
	if err != nil {
		return nil, err
	}

	return result, nil
}
