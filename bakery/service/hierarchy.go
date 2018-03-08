package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/disiqueira/gotree"
	"github.com/smartrecruiters/docker-bakery/bakery/commons"
)

// Initializes new docker hierarchy.
func NewDockerHierarchy() DockerHierarchy {
	return &dockerHierarchy{imagesWithDependantsMap: make(map[string][]*DockerImage),
		analyzedImagesFlatMap:             make(map[string]*DockerTreeItem),
		analyzedImagesSlice:               make([]*DockerTreeItem, 0),
		analyzedImagesPlusExternalParents: make([]*DockerTreeItem, 0)}
}

// Implementation of the DockerHierarchy interface
type dockerHierarchy struct {
	// map where image name is a key, value is the slice of dependant images
	imagesWithDependantsMap map[string][]*DockerImage
	// map with all analyzed images where key is the image name, value is the corresponding DockerTreeItem object
	analyzedImagesFlatMap map[string]*DockerTreeItem
	// slice with analyzed DockerTreeItem objects
	analyzedImagesSlice []*DockerTreeItem
	// slice with analyzed DockerTreeItem objects + the external parents that were not directly analyzed but should be present in the tree structure
	analyzedImagesPlusExternalParents []*DockerTreeItem
}

// Analyzes the structure of the directory and effectively builds the entire hierarchy.
// Searches for the presence of `Dockerfile.template` files.
// Uses provided map with latest versions to show it in the hierarchy.
func (h *dockerHierarchy) AnalyzeStructure(rootDir string, latestVersions map[string]*semver.Version) error {
	dockerImgParser := NewDockerImageParser()

	extractDockerImagesFn := func(sourcePath string, sourceInfo os.FileInfo, err error) error {
		name := sourceInfo.Name()
		// we can skip analysis of pure dockerfile because if it does not have a template we will not
		// be able to propagate dependency updates
		if !sourceInfo.IsDir() && name == "Dockerfile.template" {
			dockerImg, err := dockerImgParser.ParseDockerfile(sourcePath)
			if err != nil {
				return err
			}
			if dockerImg != nil {
				h.AddImage(dockerImg, latestVersions[dockerImg.Name])
			}
		}

		return nil
	}

	fmt.Println("Analyzing Dockerfile.template files")
	return filepath.Walk(rootDir, extractDockerImagesFn)
}

// Adds image to the hierarchy, uses latest version information to include it in the hierarchy view.
// Used during analyzing process.
// Updates internal hierarchy structures.
func (h *dockerHierarchy) AddImage(dockerImg *DockerImage, latestVersion *semver.Version) {
	if _, exists := h.imagesWithDependantsMap[dockerImg.DependsFromShort]; !exists {
		h.imagesWithDependantsMap[dockerImg.DependsFromShort] = make([]*DockerImage, 0)
	}
	h.imagesWithDependantsMap[dockerImg.DependsFromShort] = append(h.imagesWithDependantsMap[dockerImg.DependsFromShort], dockerImg)

	item := h.buildDockerTreeItem(dockerImg, latestVersion)
	commons.Debugf("Processing %+v", item)

	h.analyzedImagesSlice = append(h.analyzedImagesSlice, item)
	h.analyzedImagesFlatMap[dockerImg.Name] = item
	h.analyzedImagesPlusExternalParents = append(h.analyzedImagesPlusExternalParents, item)
}

// Return the map of docker images where key is the image name and value is the slice of its dependant images.
func (h *dockerHierarchy) GetImagesWithDependants() map[string][]*DockerImage {
	return h.imagesWithDependantsMap
}

// Creates the first level of the hierarchy tree. External images are qualified as first level parents.
func (h *dockerHierarchy) createFirstLevelRoots() {
	for _, i := range h.analyzedImagesSlice {
		if h.isExternalParent(i.ParentId) {
			if _, exists := h.analyzedImagesFlatMap[i.ParentId]; !exists {
				h.analyzedImagesFlatMap[i.ParentId] = &DockerTreeItem{
					Id:       i.ParentId,
					ParentId: "",
					TreeItem: &gotree.GTStructure{
						Name:  i.ParentId,
						Items: make([]*gotree.GTStructure, 0)}}

				h.analyzedImagesPlusExternalParents = append(h.analyzedImagesPlusExternalParents, h.analyzedImagesFlatMap[i.ParentId])
			}
		}
	}
}

// Checks whenever image is an external parent.
// External parent is an image that appears in the hierarchy only in the `FROM` clause but it is not defined among the hierarchy itself.
func (h *dockerHierarchy) isExternalParent(parentId string) bool {
	for _, i := range h.analyzedImagesSlice {
		if i.Id == parentId {
			return false
		}
	}
	return true
}

// Builds an object representing image in the tree view.
func (h *dockerHierarchy) buildDockerTreeItem(dockerImg *DockerImage, latestVersion *semver.Version) *DockerTreeItem {
	name := dockerImg.Name
	if latestVersion != nil {
		name = fmt.Sprintf("%s (latest: %s)", dockerImg.Name, latestVersion.String())
	}
	return &DockerTreeItem{
		Id:       dockerImg.Name,
		ParentId: dockerImg.DependsFromShort,
		TreeItem: &gotree.GTStructure{
			Name:  name,
			Items: make([]*gotree.GTStructure, 0)}}
}

// Build tree view of the hierarchy.
func (h *dockerHierarchy) buildTree(root *gotree.GTStructure) {
	for _, i := range h.analyzedImagesPlusExternalParents {
		// if item has a parent append it to the children of that parent
		if i.ParentId != "" {
			myParent := h.analyzedImagesFlatMap[i.ParentId]
			myParent.TreeItem.Items = append(myParent.TreeItem.Items, i.TreeItem)
		} else {
			// otherwise it will be treated as a root item
			root.Items = append(root.Items, i.TreeItem)
		}
	}
}

// Prints hierarchy of docker images using external slightly modified (pointers added) gotree library.
func (h *dockerHierarchy) PrintImageHierarchy(rootName string) {
	var root gotree.GTStructure
	root.Name = rootName
	h.createFirstLevelRoots()
	h.buildTree(&root)
	gotree.PrintTree(root)
}
