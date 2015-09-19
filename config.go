package main

import (
	"io/ioutil"

	"github.com/codeskyblue/go-sh"
	goyaml "gopkg.in/yaml.v2"
)

type PackageConfig struct {
	Author      string   `yaml:"author"`
	Description string   `yaml:"description"`
	Includes    []string `yaml:"includes"`
	Excludes    []string `yaml:"excludes"`
	Depth       int      `yaml:"-"`
	Script      []string `yaml:"script"`
	Settings    struct {
		TargetDir string   `yaml:"targetdir"` // target dir
		Outfiles  []string `yaml:"outfiles"`
	} `yaml:"-"`
}

const RCFILE = ".gopack.yml"

var DEFAULT_SCRIPT = []string{"go get -v", "go build"}
var DefaultPcfg *PackageConfig = &PackageConfig{
	Includes: []string{"README.md", "LICENSE", "conf", "templates", "public", "static", "views"},
	Excludes: []string{"\\.git"},
	Depth:    20,
	Script:   DEFAULT_SCRIPT,
}

// parse yaml
func ReadPkgConfig(filepath string) (pcfg *PackageConfig, err error) {
	pcfg = DefaultPcfg
	if sh.Test("file", filepath) {
		data, er := ioutil.ReadFile(filepath)
		if er != nil {
			err = er
			return
		}
		if err = goyaml.Unmarshal(data, &pcfg); err != nil {
			return
		}
	}
	return pcfg, nil
}
