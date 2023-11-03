package model

type Deployment map[string][]string

// Get contract names in deployment
func (d Deployment) ContractNames(account string) []string {
	return d[account]
}
