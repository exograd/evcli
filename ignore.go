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

		if err := is.addPattern(line); err != nil {
			return fmt.Errorf("invalid ignore entry %q: %w", line, err)
		}
	}

	return nil
}

func (is *IgnoreSet) addPattern(s string) error {
	if s[0] != '/' && !strings.HasPrefix(s, "**/") {
		// "foo/bar" matches "bar" files in a "foo" directory at any depth
		// level. We can expand it to "**/foo/bar.
		s = "**/" + s
	}

	if s[len(s)-1] == '/' && !strings.HasSuffix(s, "/**/") {
		// "/foo/bar/" means "recursively match all files in the /foo/bar/
		// directory". Therefore the glob pattern is "/foo/bar/**".
		s += "**"
	}

	glob, err := glob.Compile(s, '/')
	if err != nil {
		return fmt.Errorf("invalid glob pattern %q: %w", s, err)
	}

	is.Entries = append(is.Entries, IgnoreEntryMatch(glob))

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
