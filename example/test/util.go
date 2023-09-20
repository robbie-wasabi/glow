package test

import "fmt"

func TxPath(filename string) string {
	return fmt.Sprintf("%s/%s.cdc", "/transaction", filename)
}

func ScPath(filename string) string {
	return fmt.Sprintf("%s/%s.cdc", "/script", filename)
}

func ContractPath(filename string) string {
	return fmt.Sprintf("%s/%s.cdc", "/contract", filename)
}
