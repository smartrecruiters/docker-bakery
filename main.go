// Package main is for initializing CLI application
package main

import (
	"os"

	"fmt"

	"github.com/smartrecruiters/docker-bakery/bakery"
	cliCommons "github.com/smartrecruiters/docker-bakery/bakery/commons/cli"
	"github.com/urfave/cli"
)

const (
	applicationName        = "docker-bakery"
	applicationDescription = "CLI application for easier management of docker files"
)

var (
	version string
	commit  string
	date    string
)

func versionString() string {
	return fmt.Sprintf("%s, commit %s, built at %s", version, commit, date)
}

// Main file for the docker-bakery, assembles the app and its CLI commands
func main() {
	app := cli.NewApp()
	app.Name = applicationName
	app.Usage = applicationDescription
	app.Version = versionString()
	app.Commands = bakery.GetCommands()
	app.ExitErrHandler = cliCommons.CustomExitHandler
	app.Run(os.Args)
}
