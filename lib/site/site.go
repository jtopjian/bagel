package site

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Site represents an site file.
type Site struct {
	Roles       map[string]Role       `yaml:"roles"`
	Inventories map[string]Inventory  `yaml:"inventories"`
	Connections map[string]Connection `yaml:"connections"`
}

type Role struct {
	Inventories []string `yaml:"inventories"`
}

// New will create an Site from an site.yaml file.
func New(path string) (*Site, error) {
	site, err := readSite(path)
	if err != nil {
		return nil, err
	}

	return site, err
}

// readFile will read an site file and parse it as an Site.
func readSite(path string) (*Site, error) {
	var site Site

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s not found", path)
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading site file %s: %s", path, err)
	}

	err = yaml.Unmarshal(yamlFile, &site)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML in %s: %s", path, err)
	}

	return &site, nil
}
