package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func copyDockerImgObjectForCompare(dockerImg *DockerImage) *DockerImage {
	return &DockerImage{
		Name:             dockerImg.Name,
		DockerfilePath:   "dummy",
		DockerfileDir:    dockerImg.DockerfileDir,
		DependsOnLong:    dockerImg.DependsOnLong,
		DependsOnShort:   dockerImg.DependsOnShort,
		DependsOnVersion: dockerImg.DependsOnVersion,
		nextVersion:      dockerImg.nextVersion,
		latestVersion:    dockerImg.latestVersion,
	}
}

func TestParserDependenciesFromFromClause(t *testing.T) {
	dockerImgParser := NewDockerImageParser()

	dockerImgUnnamed, _ := dockerImgParser.ParseDockerfile("testcases/Dockerfile1")
	dockerImgUnnamedCopy := copyDockerImgObjectForCompare(dockerImgUnnamed)

	dockerImgNamedUppercase, _ := dockerImgParser.ParseDockerfile("testcases/Dockerfile2")
	dockerImgNamedUppercaseCopy := copyDockerImgObjectForCompare(dockerImgNamedUppercase)

	dockerImgNamedLowercase, _ := dockerImgParser.ParseDockerfile("testcases/Dockerfile3")
	dockerImgNamedLowercaseCopy := copyDockerImgObjectForCompare(dockerImgNamedLowercase)

	assert.Equal(t, dockerImgUnnamed.DependsOnLong, "ubuntu:2404")
	assert.Equal(t, dockerImgUnnamed.DependsOnShort, "ubuntu")
	assert.Equal(t, dockerImgUnnamed.DependsOnVersion, "2404")

	assert.Equal(t, dockerImgUnnamedCopy, dockerImgNamedUppercaseCopy)
	assert.Equal(t, dockerImgUnnamedCopy, dockerImgNamedLowercaseCopy)
	assert.Equal(t, dockerImgNamedUppercaseCopy, dockerImgNamedLowercaseCopy)
}
