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

var DefaultPcfg *PackageConfig

var DEFAULT_SCRIPT = []string{"go get -v", "go install"}

func init() {
	pcfg := &PackageConfig{}
	pcfg.Author = ""
	pcfg.Includes = []string{"README.md", "LICENSE", "conf", "templates", "public", "static", "views"}
	pcfg.Excludes = []string{"\\.git", ".*\\.go"}
	pcfg.Depth = 20
	pcfg.Settings.TargetDir = ""
	pcfg.Script = DEFAULT_SCRIPT
	DefaultPcfg = pcfg
}

// parse yaml
func ReadPkgConfig(filepath string) (pcfg PackageConfig, err error) {
	pcfg = PackageConfig{}
	if sh.Test("file", filepath) {
		data, er := ioutil.ReadFile(filepath)
		if er != nil {
			err = er
			return
		}
		if err = goyaml.Unmarshal(data, &pcfg); err != nil {
			return
		}
	} else {
		pcfg = *DefaultPcfg
	}
	return pcfg, nil
}
