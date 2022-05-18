package client

import (
	"github.com/onflow/cadence"
)

// todo: this struct might be overkill but I'm fond of the homogenous pattern it adheres to
type Sc struct {
	Script string
	// todo: any more properties?
}

func (c GlowClient) sc(cdc string) string {
	return c.replaceImportAddresses(cdc)
}

// Create new Script
func (c GlowClient) NewSc(cdc string) Sc {
	return Sc{
		Script: c.sc(cdc),
	}
}

// Create new Script from file
func (c GlowClient) NewScFromFile(file string) (*Sc, error) {
	cdc, err := c.CadenceFromFile(file)
	if err != nil {
		return nil, err
	}

	return &Sc{
		Script: c.sc(cdc),
	}, nil
}

// Execute a Script from a string
func (c GlowClient) ExecSc(
	cdc string,
	args ...cadence.Value,
) (cadence.Value, error) {
	result, err := c.Services.Scripts.Execute([]byte(cdc), args, "", c.network)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Execute a Script at a specified file path
func (c GlowClient) ExecScFromFile(
	file string,
	args ...cadence.Value,
) (cadence.Value, error) {
	sc, err := c.NewScFromFile(file)
	if err != nil {
		return nil, err
	}

	res, err := c.ExecSc(sc.Script, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
