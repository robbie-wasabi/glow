package fcl

import (
	"fmt"
)

const (
	path_base        = "/example"
	path_flow_json   = path_base + "/flow.json"
	path_transaction = "/transaction"
	path_script      = "/script"
)

func txPath(filename string) string {
	return fmt.Sprintf("%s%s/%s.cdc", path_base, path_transaction, filename)
}

func scPath(filename string) string {
	return fmt.Sprintf("%s%s/%s.cdc", path_base, path_script, filename)
}
