package service

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/smartrecruiters/go-tools/commons"
)

type postPushListener struct{}

func (pcl *postPushListener) OnPostCommand(result *CommandResult) {
	TagVersion(result.Name, result.NextVersion)
}

func NewPostPushListener() PostCommandListener {
	return &postPushListener{}
}

// Returns latest semantic version for the given image name available in the provided directory
// Version is determined based on the last git tag available for that image name
// In case no previous tags has been found the 0.0.0 is returned
func GetLatestImageVersion(imageName string) *semver.Version {
	listTagsCmd := exec.Command("git", "tag", "--list", fmt.Sprintf("%s@*", imageName))
	listTagsCmd.Dir = config.RootDir
	fmt.Printf("Executing: %s\n", listTagsCmd.Args)
	out, err := listTagsCmd.Output()
	if err != nil {
		return nil
	}

	versions := make([]*semver.Version, 0)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		ver, err := semver.NewVersion(strings.TrimPrefix(line, imageName+"@"))
		if err != nil {
			commons.Debugf("Error parsing version: %s for line: %s", err, line)
		}
		if ver != nil {
			versions = append(versions, ver)
		}
	}

	sort.Sort(semver.Collection(versions))
	if len(versions) <= 0 {
		startingVersion, _ := semver.NewVersion("0.0.0")
		return startingVersion
	}

	return versions[len(versions)-1]
}

func GetLatestVersions() (map[string]*semver.Version, error) {
	listTagsCmd := exec.Command("git", "tag", "--list")
	listTagsCmd.Dir = config.RootDir
	fmt.Println("Obtaining image latest versions from git tags")
	out, err := listTagsCmd.Output()
	if err != nil {
		return nil, err
	}

	versions := make(map[string]*semver.Version, 0)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		tag := scanner.Text()
		lineParts := strings.Split(tag, "@")
		if len(lineParts) != 2 {
			fmt.Printf("Skipping version extraction for tag: %s\n", tag)
			continue
		}
		imgName := lineParts[0]
		ver, err := semver.NewVersion(lineParts[1])
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

	return versions, nil
}

func PushTags() error {
	pushTagsCmd := exec.Command("git", "push", "--tags")
	pushTagsCmd.Dir = config.RootDir
	pushTagsCmd.Stdin = os.Stdin
	pushTagsCmd.Stdout = os.Stdout
	pushTagsCmd.Stderr = os.Stderr
	fmt.Printf("Executing: %s\n", pushTagsCmd.Args)
	return pushTagsCmd.Run()
}

func TagVersion(imageName, version string) error {
	tagCmd := exec.Command("git", "tag", fmt.Sprintf("%s@%s", imageName, version))
	tagCmd.Dir = config.RootDir
	tagCmd.Stdin = os.Stdin
	tagCmd.Stdout = os.Stdout
	tagCmd.Stderr = os.Stderr
	fmt.Printf("Executing: %s\n", tagCmd.Args)
	return tagCmd.Run()
}

// Returns the next semantic version according to the scope (major/minor/patch)
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
