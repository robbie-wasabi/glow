package client

import (
	"fmt"

	. "github.com/rrossilli/glow/util"
)

// Deployment for current network
func (c GlowClient) Deployment() map[string][]string {
	deployment := c.FlowJSON.GetDeployment(c.network)
	if IsEmpty(deployment) {
		panic(fmt.Sprintf("deployment not found in flow.json: %s", c.network))
	}
	return deployment
}

// Get deployment for current network and account
func (c GlowClient) GetAccountDeployment(name string) []string {
	contractNames := c.FlowJSON.GetAccountDeployment(c.network, name)
	return contractNames
}
