package service

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"os/signal"
	"syscall"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/smartrecruiters/docker-bakery/bakery/commons"
)

const (
	dependencyPrefix = "FROM "
	outputSeparator  = "====================================================================\n"
)

var config *Config
var dependencies map[string][]*DockerImage
var errors = make([]error, 0)
var commandResults = make([]*CommandResult, 0)
var dockerImgParser = NewDockerImageParser()
var versions map[string]*semver.Version

// Called before execution of other commands, parses config and gathers docker image dependencies/hierarchy
func InitConfiguration(configFile, rootDir string) error {
	var err error
	if configFile == "" {
		return fmt.Errorf("config file path has to be provided")
	}

	config, err = ReadConfig(configFile)
	if err != nil {
		return fmt.Errorf("could not read config file from %s due to: %s", configFile, err)
	}

	if config.RootDir == "" {
		config.RootDir = path.Dir(configFile)
		fmt.Printf("RootDir not defined in config, applying config parent dir: %s\n", config.RootDir)
	}

	if rootDir != "" {
		fmt.Printf("Overriding config rootDir to: %s\n", rootDir)
		config.RootDir = rootDir
	}

	versions, err = GetLatestVersions()
	if err != nil {
		return err
	}

	// update config properties with the latest versions of available images
	// versions are determined from the git tags
	// this is especially useful when for the first time new child image is about to be build (without being triggered by a parent build)
	// and there is no PARENT_VERSION property defined in the config file
	// in such case it will be determined dynamically from the git tags and and may be used in the child Dockerfile.template
	config.UpdateVersionProperties(versions)

	hierarchy := NewDockerHierarchy()
	err = hierarchy.AnalyzeStructure(config.RootDir, versions)
	if err != nil {
		return err
	}

	rootName := fmt.Sprintf("Dockerfiles hierarchy discovered in %s", config.RootDir)
	hierarchy.PrintImageHierarchy(rootName)
	dependencies = hierarchy.GetImagesWithDependants()

	return nil
}

// Takes the input Dockerfile.template and fills it to deliver Dockerfile that will be used to build the image.
// Uses properties defined in the config file + dynamic properties for filling the template.
// Dynamic properties are prepared automatically after analysing entire image hierarchy.
func FillTemplate(inputFile, outputFile string) error {
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		if strings.HasSuffix(inputFile, ".template") {
			return fmt.Errorf("%s does not exists", inputFile)
		}

		inputFileDir, err := dockerImgParser.ExtractDockerFileDir(inputFile)
		if err != nil {
			return err
		}

		// if inputFile does not exist check for existence of the template
		templateFile := path.Join(inputFileDir, "Dockerfile.template")
		if _, err := os.Stat(templateFile); os.IsNotExist(err) {
			return fmt.Errorf("neither %s nor %s does not exists", inputFile, templateFile)
		}

		inputFile = templateFile
	}

	if inputFile == outputFile {
		fmt.Printf("Skipping templating for %s (input path is the same as output)\n", inputFile)
		return nil
	}

	fmt.Printf("Templating %s to %s\n", inputFile, outputFile)
	return commons.FillTemplate(inputFile, outputFile, config.Properties)
}

// Uses build command defined in the config to build provided dockerfile and potentially its dependants.
// Prints the build report at the end of processing.
func BuildDockerfile(dockerfile, scope string, shouldTriggerDependantBuilds bool) error {
	defer PrintReport()
	setupInterruptionSignalHandler()
	err := ExecuteDockerCommand(config.Commands.DefaultBuildCommand, dockerfile, scope, nil, shouldTriggerDependantBuilds)
	if err != nil {
		storeError(fmt.Errorf("error processing %s: %s", dockerfile, err))
	}
	return err
}

// Uses push command defined in the config to build provided dockerfile and potentially its dependants.
// Prints the build report at the end of processing.
func PushDockerImages(dockerfile, scope string, shouldTriggerDependantBuilds bool) error {
	defer PrintReport()
	setupInterruptionSignalHandler()
	err := ExecuteDockerCommand(config.Commands.DefaultPushCommand, dockerfile, scope, NewPostPushListener(), shouldTriggerDependantBuilds)
	if err != nil {
		storeError(fmt.Errorf("error processing %s: %s", dockerfile, err))
	} else {
		err = PushTags()
	}
	return err
}

// Prints the report with processed images and its versions.
func PrintReport() {
	fmt.Printf(outputSeparator)
	fmt.Printf("Processed %d image(s):\n", len(commandResults))
	for _, r := range commandResults {
		fmt.Printf(color.GreenString("\t%s %s => %s\n", r.Name, r.CurrentVersion, r.NextVersion))
	}
	errorCount := len(errors)
	if errorCount > 0 {
		fmt.Printf(color.RedString("Following (%d) errors occurred during image processing:\n", errorCount))
		for _, err := range errors {
			fmt.Printf(color.RedString("\t%s\n", err))
		}
	}
}

// Setups interruption signal handler, that allows for printing the summary report even in cases when
// processing was aborted.
func setupInterruptionSignalHandler() {
	// setup signal receiving channel
	signalReceiverChannel := make(chan os.Signal, 3)

	// bind channel to handle signals
	signal.Notify(signalReceiverChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	// listen for channel notifications
	go func() {
		sig := <-signalReceiverChannel
		commons.Debugf("Received signal: %s\n", sig.String())
		PrintReport()
		os.Exit(1)
	}()
}

// Build/Push docker file in the following steps:
// - obtain current image info (name, version, dependants)
// - get next version based on git tags according to the change scope
// - updates dynamic config properties based on gathered info
// - templates docker command
// - execute already filled template of the build/push command
// - depending on the shouldTriggerDependantBuilds flag executes child builds if there are any
func ExecuteDockerCommand(command, dockerfile, scope string, postCmdListener PostCommandListener, shouldTriggerDependantBuilds bool) error {
	fmt.Printf(outputSeparator)
	imgName, err := dockerImgParser.ExtractImageName(dockerfile)
	if err != nil {
		return err
	}

	dockerfileDir, err := dockerImgParser.ExtractDockerFileDir(dockerfile)
	if err != nil {
		return err
	}

	latestVersion := GetLatestImageVersion(versions, imgName)
	latestVersionStr := latestVersion.String()
	nextVersion := GetNextVersion(latestVersion, scope)
	nextVersionStr := nextVersion.String()

	fmt.Printf("Working with %s scope of: %s version: %s => %s\n", scope, imgName, latestVersionStr, nextVersionStr)

	// since now we know the image name and the next version so we can
	// update config properties so that commands and dockerfile template could be properly filled
	config.UpdateDynamicProperties(imgName, nextVersionStr, dockerfileDir)
	if config.Verbose {
		config.PrintProperties()
	}

	outputPath := fmt.Sprintf("%s/Dockerfile", dockerfileDir)
	err = FillTemplate(dockerfile, outputPath)
	if err != nil {
		return err
	}

	err = executeCommand(command)
	if err != nil {
		return err
	}

	result := storeResult(imgName, latestVersionStr, nextVersionStr)

	// invoke post build listener if there is any
	if postCmdListener != nil {
		postCmdListener.OnPostCommand(result)
	}

	hasDependantImages := dependencies[result.Name] != nil && len(dependencies[result.Name]) > 0
	if shouldTriggerDependantBuilds && hasDependantImages {
		for _, dependant := range dependencies[result.Name] {
			if commons.Contains(config.AutoBuildExcludes, dependant.Name) {
				fmt.Printf("Skipping dependant build of %s as it is defined in the config autoBuildExcludes section\n", dependant.Name)
				continue
			}
			fmt.Printf("Triggering dependant build of %s\n", dependant.Name)
			err = ExecuteDockerCommand(command, dependant.DockerfilePath, scope, postCmdListener, true)
			if err != nil {
				storeError(fmt.Errorf("error processing %s: %s", dependant.Name, err))
			}
		}
	}

	return nil
}

// Executes command and prints to the stdout its output.
func executeCommand(command string) error {
	t, err := template.New("dockerCmd").Parse(command)
	var cmdBuf bytes.Buffer
	err = t.Execute(&cmdBuf, config.Properties)
	if err != nil {
		return err
	}

	dockerCmdString := cmdBuf.String()
	fmt.Printf("Executing: %s\n", dockerCmdString)
	dockerCmdWithArgs := strings.Split(dockerCmdString, " ")
	dockerCmd := exec.Command(dockerCmdWithArgs[0], dockerCmdWithArgs[1:]...)
	dockerCmd.Stdin = os.Stdin
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr

	return dockerCmd.Run()
}

// Stores the result of successful command processing. Receives image name and its current and next versions.
func storeResult(imgName string, latestVersionStr string, nextVersionStr string) *CommandResult {
	result := &CommandResult{
		Name:           imgName,
		CurrentVersion: latestVersionStr,
		NextVersion:    nextVersionStr}

	commandResults = append(commandResults, result)

	return result
}

// Stores the error of command processing.
func storeError(err error) {
	errors = append(errors, err)
}
