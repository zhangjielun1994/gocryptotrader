package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {
	tester := func(in string) error {
		err := Write(in, []byte("GoCryptoTrader"))
		if err != nil {
			return err
		}
		return os.Remove(in)
	}

	type testTable struct {
		InFile      string
		ErrExpected bool
		Cleanup     bool
	}

	var tests []testTable
	testFile := filepath.Join(os.TempDir(), "gcttest.txt")
	switch runtime.GOOS {
	case "windows":
		tests = []testTable{
			{InFile: "*", ErrExpected: true},
			{InFile: testFile, ErrExpected: false},
		}
	default:
		tests = []testTable{
			{InFile: "", ErrExpected: true},
			{InFile: testFile, ErrExpected: false},
		}
	}

	for x := range tests {
		err := tester(tests[x].InFile)
		if err != nil && !tests[x].ErrExpected {
			t.Errorf("Test %d failed, unexpected err %s\n", x, err)
		}
	}
}

func TestMove(t *testing.T) {
	tester := func(in, out string, write bool) error {
		if write {
			if err := ioutil.WriteFile(in, []byte("GoCryptoTrader"), 0770); err != nil {
				return err
			}
		}

		if err := Move(in, out); err != nil {
			return err
		}

		contents, err := ioutil.ReadFile(out)
		if err != nil {
			return err
		}

		if !strings.Contains(string(contents), "GoCryptoTrader") {
			return fmt.Errorf("unable to find previously written data")
		}

		return os.Remove(out)
	}

	type testTable struct {
		InFile      string
		OutFile     string
		Write       bool
		ErrExpected bool
	}

	var tests []testTable
	switch runtime.GOOS {
	case "windows":
		tests = []testTable{
			{InFile: "*", OutFile: "gct.txt", Write: true, ErrExpected: true},
			{InFile: "*", OutFile: "gct.txt", Write: false, ErrExpected: true},
			{InFile: "in.txt", OutFile: "*", Write: true, ErrExpected: true},
			{InFile: "in.txt", OutFile: "gct.txt", Write: true, ErrExpected: false},
		}
	default:
		tests = []testTable{
			{InFile: "", OutFile: "gct.txt", Write: true, ErrExpected: true},
			{InFile: "", OutFile: "gct.txt", Write: false, ErrExpected: true},
			{InFile: "in.txt", OutFile: "", Write: true, ErrExpected: true},
			{InFile: "in.txt", OutFile: "gct.txt", Write: true, ErrExpected: false},
		}
	}

	for x := range tests {
		err := tester(tests[x].InFile, tests[x].OutFile, tests[x].Write)
		if err != nil && !tests[x].ErrExpected {
			t.Errorf("Test %d failed, unexpected err %s\n", x, err)
		}
	}
}
