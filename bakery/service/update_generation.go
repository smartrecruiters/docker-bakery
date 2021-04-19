package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// GenerateImagesTree generate ancestors for a given image with a new parent image
func GenerateImagesTree(previousImageName string, recursive, skipExistingDirectories bool, nameReplacements []string) error {
	previousDockerImage := hierarchy.GetImageByName(previousImageName)
	if previousDockerImage == nil {
		return fmt.Errorf("unable to find image %s in the analyzed structure (is invocation directory correct?)", previousImageName)
	}

	replacer := createNewNameReplacer(nameReplacements)
	return generateAncestors(previousImageName, dependencies, replacer, recursive, skipExistingDirectories)
}

func createNewNameReplacer(nameReplacements []string) nameReplacer {
	replacer := newNameReplacer()
	for _, replacement := range nameReplacements {
		keyValuePair := strings.Split(replacement, propertyKeyValueSeparator)
		if len(keyValuePair) != 2 {
			fmt.Printf("Unable to parse provided replacement: %s - expecting key and value to be separated with '%s'", replacement, propertyKeyValueSeparator)
			continue
		}
		replacer = replacer.add(keyValuePair[0], keyValuePair[1])
	}
	return replacer
}

func generateAncestors(previousImageName string, dependenciesMap map[string][]*DockerImage, replacer nameReplacer, recursive bool, skipExistingDirectories bool) error {
	ancestorsToGenerate, err := getChildrenToUpdate(previousImageName, dependenciesMap, replacer, recursive)
	fmt.Println("Found ancestors to recreate:")
	for _, ancestor := range ancestorsToGenerate {
		fmt.Println(ancestor.childOriginImageName)
	}
	if err != nil {
		return err
	}
	for _, ancestor := range ancestorsToGenerate {
		err := generateChildProjects(ancestor, replacer, skipExistingDirectories)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateChildProjects(ancestor imageToGenerate, replacer nameReplacer, skipExistingDirectories bool) error {
	targetPath := ancestor.newImageDirectory
	exists, err := checkDirectoryExistence(targetPath, skipExistingDirectories)
	if err != nil {
		return err
	}

	if exists || ancestor.childOriginDirectory == targetPath { // when no explicit version in image name, just update the file inside
		err = updateInPlace(ancestor, replacer)
	} else {
		err = copyImageDirectory(ancestor.childOriginDirectory, targetPath)
	}
	if err != nil {
		return err
	}
	return updateImageBaseInformation(targetPath, ancestor.newParentImage, ancestor.originalParentImage)
}

func updateImageBaseInformation(path string, newParentImageName string, originalParentName string) error {
	dockerFilePath := filepath.Join(path, dockerFileTemplateName)
	read, err := ioutil.ReadFile(dockerFilePath)
	if err != nil {
		return err
	}

	newContents := strings.Replace(string(read), originalParentName, newParentImageName, -1)
	newContents = strings.Replace(newContents, dynamicImageVersionName(originalParentName), dynamicImageVersionName(newParentImageName), -1)

	return ioutil.WriteFile(dockerFilePath, []byte(newContents), 0)
}

func copyImageDirectory(childOriginDirectory string, targetPath string) error {
	copyCommand := exec.Command("cp", "-r", childOriginDirectory, targetPath)
	copyCommand.Stdin = os.Stdin
	copyCommand.Stdout = os.Stdout
	copyCommand.Stderr = os.Stderr
	return copyCommand.Run()

}

func simpleNameReplacer(oldName, newName string) func(string) string {
	return func(s string) string {
		return strings.ReplaceAll(s, oldName, newName)
	}
}

type nameReplacer struct {
	replacements []changeNameFn
}

func (n nameReplacer) add(from, to string) nameReplacer {
	n.replacements = append(n.replacements, simpleNameReplacer(from, to))
	from2 := strings.ToUpper(strings.Replace(from, "-", "_", -1))
	to2 := strings.ToUpper(strings.Replace(to, "-", "_", -1))
	n.replacements = append(n.replacements, simpleNameReplacer(from2, to2))
	return n
}

func (n nameReplacer) replace(name string) string {
	result := name
	for _, fn := range n.replacements {
		result = fn(result)
	}
	return result
}

func newNameReplacer() nameReplacer {
	return nameReplacer{replacements: make([]changeNameFn, 0)}
}

func checkDirectoryExistence(targetPath string, skipExistingDirectories bool) (bool, error) {
	if _, err := os.Stat(targetPath); err == nil {
		if skipExistingDirectories {
			fmt.Printf("directory exists: %s, ignored explicitly\n", targetPath)
			return true, nil
		}
		return true, fmt.Errorf("error - target path %s already exist", targetPath)
	}
	return false, nil
}

func getChildrenToUpdate(originImage string, dependenciesMap map[string][]*DockerImage, nameReplacer nameReplacer, recursive bool) ([]imageToGenerate, error) {
	if _, ok := dependenciesMap[originImage]; !ok {
		return nil, nil
	}

	var children []imageToGenerate
	for _, childImage := range dependenciesMap[originImage] {
		children = append(children, imageToGenerate{
			childOriginImageName: childImage.Name,
			childOriginDirectory: childImage.DockerfileDir,
			newImageDirectory:    nameReplacer.replace(childImage.DockerfileDir),
			originalParentImage:  originImage,
			newParentImage:       nameReplacer.replace(originImage),
		})
	}

	if recursive {
		for _, ch := range children {
			gatherChildren, err := getChildrenToUpdate(ch.childOriginImageName, dependenciesMap, nameReplacer, recursive)
			if err != nil {
				return nil, err
			}
			children = append(children, gatherChildren...)
		}
	}
	return children, nil
}

func updateInPlace(ancestor imageToGenerate, replacer nameReplacer) error {
	templatePath := path.Join(ancestor.newImageDirectory, dockerFileTemplateName)
	content, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return err
	}
	newContent := replacer.replace(string(content))
	return ioutil.WriteFile(templatePath, []byte(newContent), 0644)
}

type changeNameFn func(originalName string) string

type imageToGenerate struct {
	childOriginImageName string
	childOriginDirectory string
	newImageDirectory    string
	originalParentImage  string
	newParentImage       string
}

func (b imageToGenerate) String() string {
	return fmt.Sprintf("imageToGenerate{from: %s to %s}\n", b.childOriginDirectory, b.newImageDirectory)
}
