package service

import (
	"github.com/Masterminds/semver"
	"github.com/disiqueira/gotree"
)

// Config corresponds to the config structure in json file
type Config struct {
	Properties        map[string]string `json:"properties"`
	Commands          Commands          `json:"commands"`
	RootDir           string            `json:"rootDir"`
	Verbose           bool              `json:"verbose"`
	AutoBuildExcludes []string          `json:"autoBuildExcludes"`
}

// Commands is used as part of the config to contain template of build and push commands
type Commands struct {
	DefaultBuildCommand string `json:"defaultBuildCommand"`
	DefaultPushCommand  string `json:"defaultPushCommand"`
}

// DockerImage represents docker image with its parent
type DockerImage struct {
	Name             string
	DockerfileDir    string
	DockerfilePath   string
	DependsOnLong    string
	DependsOnShort   string
	DependsOnVersion string
	nextVersion      semver.Version
	latestVersion    *semver.Version
}

// PostCommandListener is an interface that allows to plugin just after docker command is executed and before any commands on children are executed
type PostCommandListener interface {
	OnPostCommand(result *CommandResult)
}

// DockerImageParser provides functionality related to parsing docker files
type DockerImageParser interface {
	ParseDockerfile(string) (*DockerImage, error)
	ExtractDockerFileDir(string) (string, error)
	ExtractImageName(string) (string, error)
}

// DockerTreeItem used in graphical representation of the image hierarchy
type DockerTreeItem struct {
	ID       string
	ParentID string
	TreeItem *gotree.GTStructure
}

// CommandResult is an outcome of the docker command
type CommandResult struct {
	Name           string
	DockerfileDir  string
	NextVersion    string
	CurrentVersion string
}

// DockerHierarchy represents hierarchy of docker images
type DockerHierarchy interface {
	// Analyzes docker files structure under given directory and constructs entire hierarchy
	AnalyzeStructure(string, map[string]*semver.Version) error
	// Adds docker image to the hierarchy based on the docker image parent
	AddImage(dockerImg *DockerImage)
	// GetImageByName returns docker image by its name. Image can be obtained after entire hierarchy has been analyzed
	GetImageByName(imageName string) *DockerImage
	// Returns map with docker images where key is the short docker image name and
	// the value is a slice of dependent images
	GetImagesWithDependants() map[string][]*DockerImage
	// Returns map with docker images where key is the short docker image name and
	// the value is docker image object
	GetImages() map[string]*DockerImage
	// Prints gathered hierarchy under a given root name
	PrintImageHierarchy(string)
}
