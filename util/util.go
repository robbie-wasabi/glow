package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"unicode/utf8"
)

// RemoveFirstChar returns a string without its first character.
func RemoveFirstChar(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

// PPrint pretty-prints any value as indented JSON.
func PPrint(v interface{}) {
	j, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}

// RemoveHexPrefix returns s without a leading "0x".
func RemoveHexPrefix(s string) string {
	if len(s) > 1 && s[:2] == "0x" {
		return s[2:]
	}
	return s
}

// PrependHexPrefix returns s prefixed with "0x" if it's not already.
func PrependHexPrefix(s string) string {
	if len(s) > 1 && s[:2] == "0x" {
		return s
	}
	return "0x" + s
}

// IsEmpty returns true if i's value is a zero value.
func IsEmpty(i interface{}) bool {
	return reflect.ValueOf(i).IsZero()
}

// Mapkey returns the first key in m whose value matches value.
func Mapkey(m map[string]interface{}, value interface{}) (string, bool) {
	for k, v := range m {
		if v == value {
			return k, true
		}
	}
	return "", false
}

// SortStringsByCharacterLength sorts arr by string length.
// If asc is true, results are ascending; otherwise descending.
func SortStringsByCharacterLength(arr []string, asc bool) []string {
	sort.SliceStable(arr, func(i, j int) bool {
		if asc {
			return len(arr[i]) < len(arr[j])
		}
		return len(arr[i]) > len(arr[j])
	})
	return arr
}

// MapKeys returns all keys from map m.
func MapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
