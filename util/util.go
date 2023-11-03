package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
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

// func (f FlowJSON) ContractNamesSortedByLength(asc bool) []string {
// 	keys := make([]string, 0, len(f.data.Contracts))
// 	for k := range f.data.Contracts {
// 		keys = append(keys, k)
// 	}
// 	sort.SliceStable(keys, func(i, j int) bool {
// 		if asc {
// 			return len(keys[i]) < len(keys[j])
// 		} else {
// 			return len(keys[i]) > len(keys[j])
// 		}
// 	})
// 	return keys
// }

// sort strings by length
func SortStringsByCharacterLength(arr []string, asc bool) []string {
	sort.SliceStable(arr, func(i, j int) bool {
		if asc {
			return len(arr[i]) < len(arr[j])
		}
		return len(arr[i]) > len(arr[j])
	})
	return arr
}

// get keys from map
func MapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
