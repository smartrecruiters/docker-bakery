package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/smartrecruiters/docker-bakery/bakery/commons"
	"github.com/smartrecruiters/gotree"
)

const dockerFileTemplateName = "Dockerfile.template"

// NewDockerHierarchy initializes new docker hierarchy.
func NewDockerHierarchy() DockerHierarchy {
	return &dockerHierarchy{
		images:                        make(map[string]*DockerImage),
		imagesWithDependantsMap:       make(map[string][]*DockerImage),
		imagesTree:                    make(map[string]*DockerTreeItem),
		imagesTreeSlice:               make([]*DockerTreeItem, 0),
		imagesTreePlusExternalParents: make([]*DockerTreeItem, 0)}
}

// Implementation of the DockerHierarchy interface
type dockerHierarchy struct {
	// map where image name is a key, value is the slice of dependant images
	imagesWithDependantsMap map[string][]*DockerImage
	// map with all analyzed images where key is the image name, value is the corresponding DockerTreeItem object
	imagesTree map[string]*DockerTreeItem
	// map with all analyzed images where key is the image name, value is the corresponding DockerImage object
	images map[string]*DockerImage
	// slice with analyzed DockerTreeItem objects
	imagesTreeSlice []*DockerTreeItem
	// slice with analyzed DockerTreeItem objects + the external parents that were not directly analyzed but should be present in the tree structure
	imagesTreePlusExternalParents []*DockerTreeItem
}

func (h *dockerHierarchy) GetImageByName(imageName string) *DockerImage {
	return h.images[imageName]
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
		if !sourceInfo.IsDir() && name == dockerFileTemplateName {
			dockerImg, err := dockerImgParser.ParseDockerfile(sourcePath)
			if err != nil {
				return err
			}
			if dockerImg != nil {
				dockerImg.latestVersion = latestVersions[dockerImg.Name]
				commons.Debugf("Adding image to hierarchy: %+v", dockerImg)
				h.AddImage(dockerImg)
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
func (h *dockerHierarchy) AddImage(dockerImg *DockerImage) {
	if _, exists := h.imagesWithDependantsMap[dockerImg.DependsOnShort]; !exists {
		h.imagesWithDependantsMap[dockerImg.DependsOnShort] = make([]*DockerImage, 0)
	}
	h.imagesWithDependantsMap[dockerImg.DependsOnShort] = append(h.imagesWithDependantsMap[dockerImg.DependsOnShort], dockerImg)

	item := h.buildDockerTreeItem(dockerImg)
	commons.Debugf("Processing %+v", item)

	h.imagesTreeSlice = append(h.imagesTreeSlice, item)
	h.imagesTree[dockerImg.Name] = item
	h.images[dockerImg.Name] = dockerImg
	h.imagesTreePlusExternalParents = append(h.imagesTreePlusExternalParents, item)
}

// Return the map of docker images where key is the image name and value is the slice of its dependant images.
func (h *dockerHierarchy) GetImagesWithDependants() map[string][]*DockerImage {
	return h.imagesWithDependantsMap
}

// Return the map of docker images where key is the image name and value is the docker image object.
func (h *dockerHierarchy) GetImages() map[string]*DockerImage {
	return h.images
}

// Creates the first level of the hierarchy tree. External images are qualified as first level parents.
func (h *dockerHierarchy) createFirstLevelRoots() {
	for _, i := range h.imagesTreeSlice {
		if h.isExternalParent(i.ParentID) {
			if _, exists := h.imagesTree[i.ParentID]; !exists {
				h.imagesTree[i.ParentID] = &DockerTreeItem{
					ID:       i.ParentID,
					ParentID: "",
					TreeItem: &gotree.GTStructure{
						Name:  i.ParentID,
						Items: make([]*gotree.GTStructure, 0)}}

				h.imagesTreePlusExternalParents = append(h.imagesTreePlusExternalParents, h.imagesTree[i.ParentID])
			}
		}
	}
}

// Checks whenever image is an external parent.
// External parent is an image that appears in the hierarchy only in the `FROM` clause but it is not defined among the hierarchy itself.
func (h *dockerHierarchy) isExternalParent(parentID string) bool {
	for _, i := range h.imagesTreeSlice {
		if i.ID == parentID {
			return false
		}
	}
	return true
}

// Builds an object representing image in the tree view.
func (h *dockerHierarchy) buildDockerTreeItem(dockerImg *DockerImage) *DockerTreeItem {
	name := dockerImg.Name
	if dockerImg.latestVersion != nil {
		name = fmt.Sprintf("%s (latest: %s)", dockerImg.Name, dockerImg.latestVersion.String())
	}
	return &DockerTreeItem{
		ID:       dockerImg.Name,
		ParentID: dockerImg.DependsOnShort,
		TreeItem: &gotree.GTStructure{
			Name:  name,
			Items: make([]*gotree.GTStructure, 0)}}
}

// Build tree view of the hierarchy.
func (h *dockerHierarchy) buildTree(root *gotree.GTStructure) {
	for _, i := range h.imagesTreePlusExternalParents {
		// if item has a parent append it to the children of that parent
		if i.ParentID != "" {
			myParent := h.imagesTree[i.ParentID]
			myParent.TreeItem.Items = append(myParent.TreeItem.Items, i.TreeItem)
		} else {
			// otherwise it will be treated as a root item
			root.Items = append(root.Items, i.TreeItem)
		}
	}
}

// Prints hierarchy of docker images using external slightly modified (pointers added) gotree library.
func (h *dockerHierarchy) PrintImageHierarchy(rootName string) {
	root := &gotree.GTStructure{}
	root.Name = rootName
	h.createFirstLevelRoots()
	h.buildTree(root)
	gotree.PrintTree(root)
}
