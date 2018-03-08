package service

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"time"

	"github.com/Masterminds/semver"
	"github.com/smartrecruiters/docker-bakery/bakery/commons"
)

// Implementation of the PostCommandListener.
type postPushListener struct{}

// OnPostCommand executes image tagging as the PostCommand action.
func (pcl *postPushListener) OnPostCommand(result *CommandResult) {
	TagVersion(result.Name, result.NextVersion)
}

// NewPostPushListener initializes new PostPushListener.
func NewPostPushListener() PostCommandListener {
	return &postPushListener{}
}

// GetLatestImageVersion returns latest semantic version for the given image name available in the provided directory
// Version is determined based on the last git tag available for that image name
// In case no previous tags has been found the 0.0.0 is returned
func GetLatestImageVersion(versions map[string]*semver.Version, imageName string) *semver.Version {
	v := versions[imageName]
	if v == nil {
		startingVersion, _ := semver.NewVersion("0.0.0")
		return startingVersion
	}

	return v
}

// GetLatestVersions returns map with latest versions of the images based on git remote tags.
// Image name is the key and latest version is the value.
func GetLatestVersions() (map[string]*semver.Version, error) {
	// we could use faster local tags to check the versions but checking the remote ones is safer in terms of version conflicts
	start := time.Now()
	listRemoteTagsCmd := exec.Command("git", "ls-remote", "--tags", "origin")
	listRemoteTagsCmd.Dir = config.RootDir
	fmt.Println("Obtaining image latest versions from git remote tags")
	out, err := listRemoteTagsCmd.Output()
	if err != nil {
		return nil, err
	}

	versions := make(map[string]*semver.Version, 0)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		refLine := scanner.Text()
		lineParts := strings.Split(refLine, "/")
		tag := lineParts[len(lineParts)-1]
		tagParts := strings.Split(tag, "@")
		if len(tagParts) != 2 {
			fmt.Printf("Skipping version extraction for tag: %s\n", tag)
			continue
		}

		imgName := tagParts[0]
		version := tagParts[1]

		ver, err := semver.NewVersion(version)
		if err != nil {
			return nil, fmt.Errorf("error parsing version: %s for tag: %s", err, tag)
		}
		if ver != nil {
			if _, versionExists := versions[imgName]; !versionExists {
				versions[imgName] = ver
			}
			if versions[imgName].LessThan(ver) {
				versions[imgName] = ver
			}
		}
	}
	commons.Debugf("Checking remote tags took: %v", time.Since(start))

	return versions, nil
}

// PushTags pushes git tags to the remote.
func PushTags() error {
	pushTagsCmd := exec.Command("git", "push", "--tags")
	pushTagsCmd.Dir = config.RootDir
	pushTagsCmd.Stdin = os.Stdin
	pushTagsCmd.Stdout = os.Stdout
	pushTagsCmd.Stderr = os.Stderr
	fmt.Printf("Executing: %s\n", pushTagsCmd.Args)
	return pushTagsCmd.Run()
}

// TagVersion creates new tag for the image with the given version.
func TagVersion(imageName, version string) error {
	tagCmd := exec.Command("git", "tag", fmt.Sprintf("%s@%s", imageName, version))
	tagCmd.Dir = config.RootDir
	tagCmd.Stdin = os.Stdin
	tagCmd.Stdout = os.Stdout
	tagCmd.Stderr = os.Stderr
	fmt.Printf("Executing: %s\n", tagCmd.Args)
	return tagCmd.Run()
}

// GetNextVersion returns the next semantic version according to the scope (major/minor/patch)
func GetNextVersion(version *semver.Version, scope string) semver.Version {
	switch scope {
	case "major":
		return version.IncMajor()
	case "minor":
		return version.IncMinor()
	case "patch":
		return version.IncPatch()
	default:
		return version.IncPatch()
	}
}
