package main

import (
	"os"

	"github.com/smartrecruiters/docker-bakery/bakery"
	cliCommons "github.com/smartrecruiters/docker-bakery/bakery/commons/cli"
	"github.com/smartrecruiters/docker-bakery/bakery/commons/version"
	"github.com/urfave/cli"
)

const (
	applicationName        = "docker-bakery"
	applicationDescription = "CLI application for easier management of docker files"
)

func main() {
	app := cli.NewApp()
	app.Name = applicationName
	app.Usage = applicationDescription
	app.Version = version.String()
	app.Commands = bakery.GetCommands()
	app.ExitErrHandler = cliCommons.CustomExitHandler
	app.Run(os.Args)
}
