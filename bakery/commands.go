// Package bakery exposes commands available for CLI usage and binds them to internal services
package bakery

import (
	"github.com/smartrecruiters/docker-bakery/bakery/commands"
	"github.com/urfave/cli"
)

// GetCommands returns CLI commands available in the docker-bakery tool
func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "fill-template",
			Aliases: []string{"prepare", "prepare-recipe"},
			Hidden:  false,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Usage: "Required. Input Dockerfile.template.",
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Required. Output Dockerfile generated from template.",
				},
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Required. Path to config.json with properties and build commands defined.",
				},
				cli.StringFlag{
					Name:  "rootDir, rd",
					Usage: "Optional. Used to override rootDir of the dockerfiles location. Can be defined in config.json, provided in this argument or determined dynamically from the base dir of config file.",
				},
				cli.StringSliceFlag{
					Name:  "property, p",
					Usage: "Optional. Allows for providing additional multiple properties that can be used during templating. Overrides properties defined in config.json file. Expected format is: -p propertyName=propertyValue",
				},
			},
			Usage:  "Used to fill Dockerfile.template file. Values needed for template are taken from the config file and from dynamic properties provided during runtime.",
			Before: commands.InitConfiguration,
			Action: commands.FillTemplateCmd,
		},
		{
			Name:   "build",
			Hidden: false,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dockerfile, d",
					Usage: "Required. Path to dockerfile/dockerfile.template file that needs to be build.",
				},
				cli.StringFlag{
					Name:  "scope, s",
					Usage: "Required. Scope of the change used to generate the next version. Can be one of: major/minor/patch.",
				},
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Required. Path to config.json with properties and build commands defined.",
				},
				cli.StringFlag{
					Name:  "root-dir, rd",
					Usage: "Optional. Used to override rootDir of the dockerfiles location. Can be defined in config.json, provided in this argument or determined dynamically from the base dir of config file.",
				},
				cli.BoolFlag{
					Name:  "skip-dependants, sd",
					Usage: "Optional. False be default. If this flag is set build of the parent will not trigger dependant builds.",
				},
				cli.StringSliceFlag{
					Name:  "property, p",
					Usage: "Optional. Allows for providing additional multiple properties that can be used during templating. Overrides properties defined in config.json file. Expected format is: -p propertyName=propertyValue",
				},
			},
			Usage:  "Used to build next version of the images in given scope. Optionally it can skip build of dependant images.",
			Before: commands.InitConfiguration,
			Action: commands.BuildDockerfileCmd,
		},
		{
			Name:   "push",
			Hidden: false,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dockerfile, d",
					Usage: "Required. Path to the dockerfile/dockerfile.template that needs to be pushed.",
				},
				cli.StringFlag{
					Name:  "scope, s",
					Usage: "Required. Scope of the change used to generate the next version. Can be one of: major/minor/patch.",
				},
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Required. Path to config.json with properties and build commands defined.",
				},
				cli.StringFlag{
					Name:  "rootDir, rd",
					Usage: "Optional. Used to override rootDir of the dockerfiles location. Can be defined in config.json, provided in this argument or determined dynamically from the base dir of config file.",
				},
				cli.BoolFlag{
					Name:  "skip-dependants, sd",
					Usage: "Optional. False be default. If this flag is set build of the parent will not trigger dependant builds.",
				},
			},
			Usage:  "Used to push next version of the images in given scope. Optionally it can skip push of dependant images.",
			Before: commands.InitConfiguration,
			Action: commands.PushDockerImagesCmd,
		},
		{
			Name:    "show-structure",
			Aliases: []string{"ss", "show-hierarchy", "hierarchy"},
			Hidden:  false,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Required. Path to config.json with properties and build commands defined.",
				},
				cli.StringFlag{
					Name:  "rootDir, rd",
					Usage: "Optional. Used to override rootDir of the dockerfiles location. Can be defined in config.json, provided in this argument or determined dynamically from the base dir of config file.",
				},
			},
			Usage:  "Used to display hierarchy of the images",
			Action: commands.InitConfiguration,
		},
		{
			Name:    "dump-latest-versions",
			Aliases: []string{"dump"},
			Hidden:  false,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Required. Path to config.json with properties and build commands defined.",
				},
				cli.StringFlag{
					Name:  "rootDir, rd",
					Usage: "Optional. Used to override rootDir of the dockerfiles location. Can be defined in config.json, provided in this argument or determined dynamically from the base dir of config file.",
				},
				cli.StringFlag{
					Name:  "file-name, file, f",
					Usage: "Optional. File name where the result data will be stored in json format.",
					Value: "docker-images.json",
				},
				cli.StringFlag{
					Name:  "exclude-dirs, e",
					Usage: "Optional. Pattern used to exclude images located in directories that match provided argument.",
				},
			},
			Usage:  "Used to dump data about latest versions of images to the provided file",
			Before: commands.InitConfiguration,
			Action: commands.DumpLatestVersionsCmd,
		},
		{
			Name:    "copy-images-hierarchy",
			Aliases: []string{"cph"},
			Hidden:  false,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Required. Path to config.json with properties and build commands defined.",
				},
				cli.StringFlag{
					Name:  "rootDir, rd",
					Usage: "Optional. Used to override rootDir of the dockerfiles location. Can be defined in config.json, provided in this argument or determined dynamically from the base dir of config file.",
				},
				cli.StringFlag{
					Name:  "base-image, image",
					Usage: "Base image name whose hierarchy will be copied",
				},
				cli.BoolFlag{
					Name:  "recursive, r",
					Usage: "Optional. When set to true, will generate whole images family. False generates only 1 level of children",
				},
				cli.StringSliceFlag{
					Name:  "skip-existing-dirs",
					Usage: "Optional. When set to true, ignore error when target path during copying already exist",
				},
				cli.StringSliceFlag{
					Name:  "replace, rp",
					Usage: "Optional. Allows for providing additional string replacements, to automate updates in new family. Will change all occurrences in the child images. Expected format is: -rp originalString=replacementString. Example: -rp python3.8=python3.9",
				},
			},
			Usage:  "This command takes 'new-parent-image' image and copies images hierarchy from 'base-image', for example to simplify new image versions when updating some dependencies (for example from Ubuntu 18 to Ubuntu 20, without loosing previous images) ",
			Before: commands.InitConfiguration,
			Action: commands.GenerateImagesTree,
		},
	}
}
