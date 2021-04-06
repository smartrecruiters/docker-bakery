// Package commands gathers commands exposed in the CLI docker-bakery interface
package commands

import (
	"github.com/smartrecruiters/docker-bakery/bakery/service"
	"github.com/urfave/cli"
)

// InitConfiguration initializes configuration for the rest of invoked commands.
// Receives config file path and optionally root directory to override the config section.
func InitConfiguration(c *cli.Context) error {
	return service.InitConfiguration(c.String("c"), c.String("rd"), c.StringSlice("p"))
}

// FillTemplateCmd fills input dockerfile template and stores the result under provided output.
func FillTemplateCmd(c *cli.Context) error {
	return service.FillTemplate(c.String("i"), c.String("o"))
}

// BuildDockerfileCmd invokes docker build command on the provided file with the provided change scope.
// Optionally it skips builds of dependant images.
func BuildDockerfileCmd(c *cli.Context) error {
	return service.BuildDockerfile(c.String("d"), c.String("s"), !c.Bool("sd"))
}

// PushDockerImagesCmd invokes docker push command on the provided file with the provided change scope.
// Optionally it skips pushes of dependant images.
func PushDockerImagesCmd(c *cli.Context) error {
	return service.PushDockerImages(c.String("d"), c.String("s"), !c.Bool("sd"))
}

// DumpLatestVersionsCmd dumps information about images and their latest versions to file in json format.
func DumpLatestVersionsCmd(c *cli.Context) error {
	return service.DumpLatestVersions(c.String("f"), c.String("e"))
}

// GenerateImagesTree generate ancestors for a given image with a new parent image
func GenerateImagesTree(c *cli.Context) error {
	return service.GenerateImagesTree(c.String("base-image"), c.Bool("r"), c.Bool("skip-existing-dirs"), c.StringSlice("replace"))
}
