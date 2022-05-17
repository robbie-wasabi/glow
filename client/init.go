package client

import (
	"errors"
	"flag"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	MODE       = ""
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	ROOT = filepath.Join(filepath.Dir(b), "..")
)

func init() {
	var _ = func() bool {
		testing.Init()
		return true
	}()
	e := flag.String("env", "emulator", "environment to point integration test at")
	s := flag.String("mode", "", "run tests in short mode")
	// fmt.Println(s)
	flag.Parse()
	if e != nil {
		MODE = *s
	} else {
		panic(errors.New("env flag was not parsed"))
	}
}
