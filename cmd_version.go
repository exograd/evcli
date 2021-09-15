package main

import (
	"fmt"

	"github.com/galdor/go-cmdline"
)

func cmdVersion(args []string, app *App) {
	cl := cmdline.New()
	cl.Parse(args)

	fmt.Printf("evcli %s\n", buildId)
}
