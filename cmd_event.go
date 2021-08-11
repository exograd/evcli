package main

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/galdor/go-cmdline"
)

func cmdEvent(args []string, app *App) {
	cl := cmdline.New()
	cl.AddCommand("create", "create a new custom event")
	cl.Parse(args)

	var cmd func([]string, *App)

	switch cl.CommandName() {
	case "create":
		cmd = cmdEventCreate
	}

	cmd(cl.CommandNameAndArguments(), app)
}

func cmdEventCreate(args []string, app *App) {
	cl := cmdline.New()
	cl.AddArgument("connector", "the name of the connector")
	cl.AddArgument("event", "the name of the event")
	cl.AddArgument("data",
		"the JSON object representing event data (\"-\" to read stdin)")
	cl.AddOption("t", "event-time", "timestamp",
		"the date and time the event occurred (RFC 3339 format)")
	cl.Parse(args)

	var eventTime time.Time
	if cl.IsOptionSet("event-time") {
		eventTimeString := cl.OptionValue("event-time")

		var err error
		eventTime, err = time.Parse(time.RFC3339, eventTimeString)
		if err != nil {
			die("invalid event time: %v", err)
		}
	} else {
		eventTime = time.Now().UTC()
	}

	connector := cl.ArgumentValue("connector")
	name := cl.ArgumentValue("event")

	data := []byte(cl.ArgumentValue("data"))
	if len(data) == 1 && data[0] == '-' {
		var err error
		data, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			die("cannot read stdin: %v", err)
		}
	}

	newEvent := NewEvent{
		EventTime: &eventTime,
		Connector: connector,
		Name:      name,
		Data:      data,
	}

	_, err := app.Client.CreateEvent(&newEvent)
	if err != nil {
		die("cannot create event: %v", err)
	}

	info("events created")
}
