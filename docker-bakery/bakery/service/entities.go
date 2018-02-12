package service

import (
	"github.com/Masterminds/semver"
	"github.com/disiqueira/gotree"
)

// Corresponds to the config structure in json file
type Config struct {
	Properties        map[string]string `json:"properties"`
	Commands          Commands          `json:"commands"`
	RootDir           string            `json:"rootDir"`
	Verbose           bool              `json:"verbose"`
	AutoBuildExcludes []string          `json:"autoBuildExcludes"`
}

// Used as part of the config to contain template of build and push commands
type Commands struct {
	DefaultBuildCommand string `json:"defaultBuildCommand"`
	DefaultPushCommand  string `json:"defaultPushCommand"`
}

// Represents docker image with its parent
type DockerImage struct {
	Name             string
	DockerfileDir    string
	DockerfilePath   string
	DependsFromLong  string
	DependsFromShort string
}

// Interface that allows to plugin just after docker command is executed and before any commands on children are executed
type PostCommandListener interface {
	OnPostCommand(result *CommandResult)
}

// Provides functionality related to parsing docker files
type DockerImageParser interface {
	ParseDockerfile(string) (*DockerImage, error)
	ExtractDockerFileDir(string) (string, error)
	ExtractImageName(string) (string, error)
}

// Used in graphical representation of the image hierarchy
type DockerTreeItem struct {
	Id       string
	ParentId string
	TreeItem *gotree.GTStructure
}

// Outcome of the docker command
type CommandResult struct {
	Name           string
	DockerfileDir  string
	NextVersion    string
	CurrentVersion string
}

// Represents hierarchy of docker images
type DockerHierarchy interface {
	// Analyzes docker files structure under given directory and constructs entire hierarchy
	AnalyzeStructure(string, map[string]*semver.Version) error
	// Adds docker image to the hierarchy based on the docker image parent
	AddImage(dockerImg *DockerImage, latestVersion *semver.Version)
	// Returns map with docker images where key is the short docker image name and
	// the value is a slice of dependent images
	GetImagesWithDependants() map[string][]*DockerImage
	// Prints gathered hierarchy under a given root name
	PrintImageHierarchy(string)
}
