package client

import (
	"context"
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-cli/flowkit"
	"github.com/onflow/flow-go-sdk"
)

// Sc struct encapsulates Flow script execution logic.
type Sc struct {
	ctx    context.Context
	script *flowkit.Script
	client *GlowClient
}

// newSc is a utility function to create a script, private to ensure a single point of instantiation.
func (c *GlowClient) newSc(content string, args ...cadence.Value) *Sc {
	b := []byte(c.replaceImportAddresses(content))
	return &Sc{
		script: &flowkit.Script{
			Code: b,
			Args: args,
		},
		client: c,
	}
}

// NewSc creates a new script from a byte array.
func (c *GlowClient) NewSc(bytes []byte, args ...cadence.Value) *Sc {
	return c.newSc(string(bytes), args...)
}

// NewScFromString creates a new script from a string.
func (c *GlowClient) NewScFromString(cdc string, args ...cadence.Value) *Sc {
	return c.newSc(cdc, args...)
}

// NewScFromFile creates a new script from a file.
func (c *GlowClient) NewScFromFile(file string, args ...cadence.Value) (*Sc, error) {
	cdc, err := c.CadenceFromFile(file)
	if err != nil {
		return nil, fmt.Errorf("sc not found at: %s", file)
	}
	return c.newSc(cdc, args...), nil
}

// WithContext adds context to a script.
func (s *Sc) WithContext(ctx context.Context) *Sc {
	s.ctx = ctx
	return s
}

// Args sets the arguments for the script.
func (sc *Sc) Args(args ...cadence.Value) *Sc {
	sc.script.Args = args
	return sc
}

// AddArg appends a single argument to the script.
func (sc *Sc) AddArg(arg cadence.Value) *Sc {
	sc.script.Args = append(sc.script.Args, arg)
	return sc
}

// Exec executes the script at the latest block.
func (sc *Sc) Exec() (cadence.Value, error) {
	query := flowkit.ScriptQuery{
		Latest: true,
		ID:     flow.EmptyID,
		Height: 0,
	}
	return sc.client.FlowKit.ExecuteScript(sc.ctx, *sc.script, query)
}

// ExecWithQuery executes the script using a specific query.
func (sc *Sc) ExecWithQuery(query *flowkit.ScriptQuery) (cadence.Value, error) {
	return sc.client.FlowKit.ExecuteScript(sc.ctx, *sc.script, *query)
}
