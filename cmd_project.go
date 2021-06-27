package main

import (
	"github.com/galdor/go-cmdline"
)

func cmdProject(args []string, client *Client) {
	cl := cmdline.New()
	cl.AddCommand("list", "list projects")
	cl.AddCommand("create", "create a project")
	cl.AddCommand("delete", "delete a project")
	cl.AddCommand("deploy", "deploy a project")
	cl.Parse(args)

	var cmd func([]string, *Client)

	switch cl.CommandName() {
	case "list":
		cmd = cmdProjectList
	case "create":
		cmd = cmdProjectCreate
	case "delete":
		cmd = cmdProjectDelete
	case "deploy":
		cmd = cmdProjectDeploy
	}

	cmd(cl.CommandNameAndArguments(), client)
}

func cmdProjectList(args []string, client *Client) {
	cl := cmdline.New()
	cl.Parse(args)

	// TODO
	die("unimplemented")
}

func cmdProjectCreate(args []string, client *Client) {
	cl := cmdline.New()
	cl.Parse(args)

	// TODO
	die("unimplemented")
}

func cmdProjectDelete(args []string, client *Client) {
	cl := cmdline.New()
	cl.Parse(args)

	// TODO
	die("unimplemented")
}

func cmdProjectDeploy(args []string, client *Client) {
	cl := cmdline.New()
	cl.Parse(args)

	// TODO
	die("unimplemented")
}
