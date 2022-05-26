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

// Create new Script
func (c *GlowClient) NewSc(bytes []byte, args ...cadence.Value) *Sc {
	b := []byte(c.replaceImportAddresses(string(bytes)))
	return &Sc{
		cdc:    b,
		args:   args,
		client: c,
	}
}

// Create new Script from string
func (c *GlowClient) NewScFromString(cdc string, args ...cadence.Value) *Sc {
	b := []byte(c.replaceImportAddresses(cdc))
	return &Sc{
		cdc:    b,
		args:   args,
		client: c,
	}
}

// Create new Script from file
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

func (sc *Sc) Args(args ...cadence.Value) *Sc {
	sc.args = args
	return sc
}

func (sc *Sc) AddArg(arg cadence.Value) *Sc {
	sc.args = append(sc.args, arg)
	return sc
}

func (sc *Sc) Exec() (cadence.Value, error) {
	result, err := sc.client.Services.Scripts.Execute(sc.cdc, sc.args, "", sc.client.network)
	if err != nil {
		return nil, err
	}

	return result, nil
}
