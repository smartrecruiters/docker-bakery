package main

import (
	"os"

	"github.com/smartrecruiters/go-tools/docker-bakery/bakery"
	"github.com/smartrecruiters/go-tools/version"
	"github.com/urfave/cli"
	"github.com/smartrecruiters/go-tools/cli-commons"
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
	app.ExitErrHandler = cli_commons.CustomExitHandler
	app.Run(os.Args)
}
