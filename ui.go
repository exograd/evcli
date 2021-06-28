package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Confirm(prompt string) bool {
	fmt.Printf("%s\n[yn] ", prompt)

	r := bufio.NewReader(os.Stdin)
	line, _, err := r.ReadLine()
	if err != nil {
		die("cannot read stdin: %v", err)
	}

	response := strings.ToLower(strings.TrimSpace(string(line)))

	switch response {
	case "y":
		fallthrough
	case "yes":
		return true
	}

	return false
}
