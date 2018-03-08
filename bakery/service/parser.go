package service

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/smartrecruiters/docker-bakery/bakery/commons"
)

// Type holding image parser functionality.
type dockerImageParser struct{}

// Extracts directory name of the provided dockerfile.
func (dip *dockerImageParser) ExtractDockerFileDir(dockerfile string) (string, error) {
	dockerfile, err := filepath.Abs(dockerfile)
	commons.Debugf("Resolved dockerfile path to %s", dockerfile)
	if err != nil {
		return "", err
	}
	return path.Dir(dockerfile), nil
}

// Extracts name of the dockerfile based on its parent dir.
func (dip *dockerImageParser) ExtractImageName(dockerfile string) (string, error) {
	dir, err := dip.ExtractDockerFileDir(dockerfile)
	if err != nil {
		return "", err
	}
	return path.Base(dir), nil
}

// Parses dockerfile and return the object describing it.
// Apart from image name and location the parent information is extracted based on the `FROM` clause.
func (dip *dockerImageParser) ParseDockerfile(dockerfilePath string) (*DockerImage, error) {
	dockerfileDir, err := dip.ExtractDockerFileDir(dockerfilePath)
	if err != nil {
		return nil, err
	}
	imageName, err := dip.ExtractImageName(dockerfilePath)
	if err != nil {
		return nil, err
	}
	commons.Debugf("Resolved image name %s, dir: %s", imageName, dockerfileDir)
	inFile, err := os.Open(dockerfilePath)
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, dependencyPrefix) {
			return nil, fmt.Errorf("unable to extract dependency from %s file. Check if first line starts with `FROM `", dockerfilePath)
		}

		dependsFromLong := strings.TrimPrefix(line, dependencyPrefix)
		dependsFromShort := dependsFromLong
		parts := strings.Split(dependsFromLong, "/")
		if len(parts) <= 0 {
			fmt.Printf("WARN: Unable to determine short base image name for: %s", dockerfilePath)
		} else {
			imgNameWithVersion := parts[len(parts)-1]
			imgNameWithVersionParts := strings.Split(imgNameWithVersion, ":")
			dependsFromShort = imgNameWithVersionParts[0]
		}

		return &DockerImage{
			Name:             imageName,
			DependsFromLong:  dependsFromLong,
			DependsFromShort: dependsFromShort,
			DockerfileDir:    dockerfileDir,
			DockerfilePath:   dockerfilePath}, nil
	}

	err = scanner.Err()
	return nil, err
}

// Initialized new docker image parser.
func NewDockerImageParser() DockerImageParser {
	return &dockerImageParser{}
}
