package model

type Contract struct {
	Source  string            `json:"source"`
	Aliases map[string]string `json:"aliases"`
}

func (c Contract) Address(env string) string {
	return c.Aliases[env]
}
