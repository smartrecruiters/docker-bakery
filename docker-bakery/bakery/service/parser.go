package service

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/smartrecruiters/go-tools/commons"
)

type dockerImageParser struct{}

func (dip *dockerImageParser) ExtractDockerFileDir(dockerfile string) (string, error) {
	dockerfile, err := filepath.Abs(dockerfile)
	commons.Debugf("Resolved dockerfile path to %s\n", dockerfile)
	if err != nil {
		return "", err
	}
	return path.Dir(dockerfile), nil
}

func (dip *dockerImageParser) ExtractImageName(dockerfile string) (string, error) {
	dir, err := dip.ExtractDockerFileDir(dockerfile)
	if err != nil {
		return "", err
	}
	return path.Base(dir), nil
}

func (dip *dockerImageParser) ParseDockerfile(dockerfilePath string) (*DockerImage, error) {
	dockerfileDir, err := dip.ExtractDockerFileDir(dockerfilePath)
	if err != nil {
		return nil, err
	}
	imageName, err := dip.ExtractImageName(dockerfilePath)
	if err != nil {
		return nil, err
	}
	commons.Debugf("Resolved image name %s, dir: %s\n", imageName, dockerfileDir)
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

func NewDockerImageParser() DockerImageParser {
	return &dockerImageParser{}
}
