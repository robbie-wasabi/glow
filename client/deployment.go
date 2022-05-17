package client

import (
	"fmt"

	. "github.com/rrossilli/glow/util"
)

// Deployment for current env
func (c GlowClient) Deployment() map[string][]string {
	deployment := c.FlowJSON.GetDeployment(c.env)
	if IsEmpty(deployment) {
		panic(fmt.Sprintf("deployment not found in flow.json: %s", c.env))
	}
	return deployment
}

// Get deployment for current env and account
func (c GlowClient) GetAccountDeployment(name string) []string {
	contractNames := c.FlowJSON.GetAccountDeployment(c.env, name)
	return contractNames
}
