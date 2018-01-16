package commands

import (
	"github.com/smartrecruiters/go-tools/docker-bakery/bakery/service"
	"github.com/urfave/cli"
)

func InitConfiguration(c *cli.Context) error {
	return service.InitConfiguration(c.String("c"), c.String("rd"))
}

func FillTemplateCmd(c *cli.Context) error {
	return service.FillTemplate(c.String("i"), c.String("o"))
}

func BuildDockerfileCmd(c *cli.Context) error {
	return service.BuildDockerfile(c.String("d"), c.String("s"), !c.Bool("sd"))
}

func PushDockerImagesCmd(c *cli.Context) error {
	return service.PushDockerImages(c.String("d"), c.String("s"), !c.Bool("sd"))
}
