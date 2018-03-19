package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"time"

	"bytes"
	"text/template"

	"github.com/Masterminds/semver"
	"github.com/smartrecruiters/docker-bakery/bakery/commons"
)

const (
	builderNamePropName    = "BAKERY_BUILDER_NAME"
	builderEmailPropName   = "BAKERY_BUILDER_EMAIL"
	builderHostPropName    = "BAKERY_BUILDER_HOST"
	buildDatePropName      = "BAKERY_BUILD_DATE"
	imageHierarchyPropName = "BAKERY_IMAGE_HIERARCHY"
	signatureValuePropName = "BAKERY_SIGNATURE_VALUE"
	signatureEnvsPropName  = "BAKERY_SIGNATURE_ENVS"
	imageVersionPropName   = "IMAGE_VERSION"
	imageNamePropName      = "IMAGE_NAME"
	dockerfileDirPropName  = "DOCKERFILE_DIR"
)

// ReadConfig reads configuration file from provided path and returns it as an object.
func ReadConfig(configFile string) (*Config, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// UpdateDynamicProperties updates config object state with the corresponding values of all dynamic properties.
// Called in every cycle of executing docker command.
func (cfg *Config) UpdateDynamicProperties(dockerImg *DockerImage) {
	nextVersion := dockerImg.GetNextVersionString()
	buildDate := time.Now().Format("2006-01-02 15:04:05")
	cfg.setBuildDate(buildDate)
	cfg.setImageName(dockerImg.Name)
	cfg.setDockerfileDir(dockerImg.DockerfileDir)
	cfg.setConstantImageVersion(nextVersion)
	cfg.setDynamicImageVersionProperty(dockerImg.Name, nextVersion)
	cfg.setImageHierarchy(dockerImg)
	cfg.buildSignature()
}

// UpdateVersionProperties updates config with versions of the images.
// Called once after latest versions of images are known (after analysing entire structure)
func (cfg *Config) UpdateVersionProperties(versions map[string]*semver.Version) {
	for image, version := range versions {
		cfg.setDynamicImageVersionProperty(image, version.String())
	}
}

// Sets the dynamic version property for the image and updates the config with it.
// Version property follows the convention comparing to the image name:
// - it is in uppercase
// - has `-` and `.` replaced with `_`
// - has the `_VERSION` suffix
// The result of invocation setDynamicImageVersionProperty('redis', '1.0.0')
// would be property REDIS_VERSION = '1.0.0'
func (cfg *Config) setDynamicImageVersionProperty(imgName, version string) {
	imgNameInTmpl := strings.ToUpper(imgName)
	imgNameInTmpl = strings.Replace(imgNameInTmpl, "-", "_", -1)
	imgNameInTmpl = strings.Replace(imgNameInTmpl, ".", "_", -1)
	propertyName := fmt.Sprintf("%s_VERSION", imgNameInTmpl)
	commons.Debugf("Setting property %s to %s", propertyName, version)
	cfg.Properties[propertyName] = version
}

// setImageVersion sets the IMAGE_VERSION property to the provided version.
// In contrast to setDynamicImageVersionProperty it always works on the same IMAGE_VERSION property
// that is needed in templating docker build/push command
func (cfg *Config) setConstantImageVersion(version string) {
	cfg.Properties[imageVersionPropName] = version
}

func (cfg *Config) setImageName(name string) {
	cfg.Properties[imageNamePropName] = name
}

func (cfg *Config) setDockerfileDir(dir string) {
	cfg.Properties[dockerfileDirPropName] = dir
}

func (cfg *Config) setBuildDate(buildDate string) {
	cfg.Properties[buildDatePropName] = buildDate
}

func (cfg *Config) setBuilderName(name string) {
	cfg.Properties[builderNamePropName] = name
}

func (cfg *Config) setBuilderEmail(email string) {
	cfg.Properties[builderEmailPropName] = email
}

func (cfg *Config) setBuilderHost(host string) {
	cfg.Properties[builderHostPropName] = host
}

// setImageHierarchy updates the config properties with a special BAKERY_IMAGE_HIERARCHY property.
// Embedding this property in the chain of docker images allows for tracking entire hierarchy of the image including its
// parents and versions
func (cfg *Config) setImageHierarchy(dockerImage *DockerImage) {
	parentVersionSuffix := ""
	parentVersion := resolveVersion(dockerImage.DependsOnVersion, cfg)
	parentImageName := dockerImage.DependsOnShort
	if len(parentVersion) > 0 {
		parentVersionSuffix = fmt.Sprintf(":%s", parentVersion)
	}
	commons.Debugf("Resolved parent version to: %s for: %s image", parentVersion, dockerImage.Name)

	cfg.Properties[imageHierarchyPropName] = fmt.Sprintf("${BAKERY_IMAGE_HIERARCHY:-\"%s%s\"}->%s:%s",
		parentImageName,
		parentVersionSuffix,
		dockerImage.Name,
		dockerImage.GetNextVersionString())
}

// resolveVersion tries to resolve version string using dynamic properties from config
// or falls back to provided version if it was not templated.
func resolveVersion(versionToResolve string, cfg *Config) string {
	if len(versionToResolve) == 0 {
		return versionToResolve
	}

	t := template.New("version-substitution-template")
	t, err := t.Parse(versionToResolve)
	if err != nil {
		commons.Debugf("Unable to parse version as a template, falling back to non templated version [%s], err: %s", versionToResolve, err)
		return versionToResolve
	}

	buf := bytes.NewBufferString("")
	t.Execute(buf, cfg.Properties)
	return buf.String()
}

// buildSignature builds docker-bakery signature from several dynamic properties
// and updates the config with new variables available for usage in templates.
func (cfg *Config) buildSignature() {
	signatureProperties := []string{builderNamePropName, builderEmailPropName, builderHostPropName, buildDatePropName, imageHierarchyPropName}
	propsCount := len(signatureProperties)

	var signatureValueBuf bytes.Buffer
	var signatureEnvsBuf bytes.Buffer
	for i, p := range signatureProperties {
		signatureValueBuf.WriteString(cfg.Properties[p])
		signatureEnvsBuf.WriteString(fmt.Sprintf("%s=\"%s\"", p, cfg.Properties[p]))
		isNotLastProperty := i < propsCount-1
		if isNotLastProperty {
			signatureValueBuf.WriteString(";")
			signatureEnvsBuf.WriteString(" \\\n")
		}
	}

	cfg.Properties[signatureValuePropName] = signatureValueBuf.String()
	cfg.Properties[signatureEnvsPropName] = signatureEnvsBuf.String()
}

// PrintProperties prints all properties available in the config (along with the dynamic ones).
func (cfg *Config) PrintProperties() {
	fmt.Println("Config properties:")
	sortedKeys := commons.SortMapKeys(cfg.Properties)
	for _, key := range sortedKeys {
		fmt.Printf("\t%s=%s\n", key, config.Properties[key])
	}
}
