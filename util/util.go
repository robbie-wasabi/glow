package client

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unicode/utf8"
)

// remove first char of string
func RemoveFirstChar(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

// pretty print a struct
func PPrint(m interface{}) {
	j, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}

// remove "0x" from beginning of string
func RemoveHexPrefix(s string) string {
	if s[:2] == "0x" {
		return s[2:]
	}
	return s
}

// prepend "0x" to beginning of string
func PrependHexPrefix(s string) string {
	if s[:2] == "0x" {
		return s
	}
	return "0x" + s
}

// check if the value of a struct is empty
func IsEmpty(i interface{}) bool {
	return reflect.ValueOf(i).IsZero()
}

// get key based on value
func Mapkey(m map[string]interface{}, value interface{}) (key string, ok bool) {
	for k, v := range m {
		if v == value {
			key = k
			ok = true
			return
		}
	}
	return
}

func TxPath(filename string) string {
	return fmt.Sprintf("%s/%s.cdc", "/transaction", filename)
}

func ScPath(filename string) string {
	return fmt.Sprintf("%s/%s.cdc", "/script", filename)
}
