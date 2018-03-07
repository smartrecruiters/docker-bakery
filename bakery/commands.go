package bakery

import (
	"github.com/smartrecruiters/docker-bakery/bakery/commands"
	"github.com/urfave/cli"
)

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
	}
}
