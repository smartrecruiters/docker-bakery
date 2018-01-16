package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/smartrecruiters/go-tools/commons"
)

// Reads configuration file from provided path and returns it as an object
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

// Updates config object state with the corresponding values of dynamic properties
func (cfg *Config) UpdateDynamicProperties(imgName, nextVersion, dockerfileDir string) {
	cfg.Properties["IMAGE_NAME"] = imgName
	cfg.Properties["IMAGE_VERSION"] = nextVersion
	cfg.Properties["DOCKERFILE_DIR"] = dockerfileDir
	cfg.setImageVersionProperty(imgName, nextVersion)
}

func (cfg *Config) UpdateImageVersions(versions map[string]*semver.Version) {
	for image, version := range versions {
		cfg.setImageVersionProperty(image, version.String())
	}
}

func (cfg *Config) setImageVersionProperty(imgName, version string) {
	imgNameInTmpl := strings.ToUpper(imgName)
	imgNameInTmpl = strings.Replace(imgNameInTmpl, "-", "_", -1)
	imgNameInTmpl = strings.Replace(imgNameInTmpl, ".", "_", -1)
	propertyName := fmt.Sprintf("%s_VERSION", imgNameInTmpl)
	commons.Debugf("Setting property %s to %s", propertyName, version)
	cfg.Properties[propertyName] = version
}

func (cfg *Config) PrintProperties() {
	fmt.Println("Config properties:")
	sortedKeys := commons.SortMapKeys(cfg.Properties)
	for _, key := range sortedKeys {
		fmt.Printf("\t%s=%s\n", key, config.Properties[key])
	}
}
