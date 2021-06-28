package main

import (
	"net/url"

	"github.com/galdor/go-cmdline"
)

func cmdAPI(args []string, app *App) {
	cl := cmdline.New()
	cl.AddCommand("status", "print the current status of the api")
	cl.Parse(args)

	var cmd func([]string, *App)

	switch cl.CommandName() {
	case "status":
		cmd = cmdAPIStatus
	}

	cmd(cl.CommandNameAndArguments(), app)
}

func cmdAPIStatus(args []string, app *App) {
	cl := cmdline.New()
	cl.Parse(args)

	var status APIStatus

	uri := &url.URL{Path: "/v0/status"}
	err := app.Client.SendRequest("GET", uri, nil, &status)
	if err != nil {
		die("cannot fetch api status: %v", err)
	}
}
