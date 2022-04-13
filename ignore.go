package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/gobwas/glob"
)

type IgnoreSet struct {
	Entries []IgnoreEntry
}

type IgnoreEntry interface{}

type IgnoreEntryMatch = glob.Glob

func (is *IgnoreSet) LoadDirectoryIfExists(dirPath string) error {
	filePath := path.Join(dirPath, ".evcli-ignore")
	return is.LoadFileIfExists(filePath)
}

func (is *IgnoreSet) LoadFileIfExists(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("NOT FOUND CACA: %s\n", filePath)
		return nil
	} else if err != nil {
		return fmt.Errorf("cannot read %s: %w", filePath, err)
	}

	p.Debug(1, "loading ignore set from %s", filePath)

	return is.LoadData(data)
}

func (is *IgnoreSet) LoadData(data []byte) error {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		glob, err := glob.Compile(line, '/')
		if err != nil {
			return fmt.Errorf("invalid glob pattern %q: %w", line, err)
		}

		entry := IgnoreEntryMatch(glob)

		is.Entries = append(is.Entries, entry)
	}

	return nil
}

func (is *IgnoreSet) Match(filePath string) bool {
	for _, e := range is.Entries {
		switch v := e.(type) {
		case IgnoreEntryMatch:
			if v.Match(filePath) {
				return true
			}

		default:
			panic(fmt.Errorf("unhandled ignore set entry of type %T", e))
		}
	}

	return false
}
