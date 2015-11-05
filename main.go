package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/codegangsta/cli"
	"github.com/gorelease/gopack/cmds"
	goyaml "gopkg.in/yaml.v2"
)

const VERSION = "0.2.0915"

var app = cli.NewApp()

func inputString(key string, defa string) string {
	fmt.Printf("[?] %s: (%s) ", key, defa)
	var val string
	if _, err := fmt.Scanf("%s", &val); err != nil || val == "" {
		return defa
	}
	return val
}

func gitUsername() string {
	output, _ := exec.Command("git", "config", "user.name").Output()
	username := string(output)
	if username == "" {
		return "unknown"
	}
	return strings.TrimSpace(username)
}

func actionInit(ctx *cli.Context) {
	if _, err := os.Stat(RCFILE); err == nil && !ctx.Bool("force") {
		fmt.Printf("config file %s already exists\n", RCFILE)
		return
	}

	pcfg := DefaultPcfg
	pcfg.Author = inputString("author", gitUsername())
	pcfg.Description = inputString("description", "...")

	data, _ := goyaml.Marshal(DefaultPcfg)
	beautiData := strings.Replace(string(data), "\n-", "\n  -", -1)
	ioutil.WriteFile(RCFILE, []byte(beautiData), 0644)
	fmt.Println("Configuration file save to .gopack.yml")
}

type OSArch struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

func actionAll(ctx *cli.Context) {
	ss := map[string][]string{
		"darwin":  {"amd64"},
		"windows": {"amd64", "386"},
		"linux":   {"amd64", "386", "arm"},
	}

	oses := strings.Fields(ctx.String("os"))
	if len(oses) == 0 {
		oses = []string{"linux", "darwin", "windows"}
	}
	pathTemplate := ctx.String("output")
	osarches := make([]OSArch, 0)
	for _, os := range oses {
		for _, arch := range ss[os] {
			osarches = append(osarches, OSArch{os, arch})
		}
	}

	for _, oa := range osarches {
		tmpl := template.Must(template.New("path").Parse(pathTemplate))
		cwd, _ := os.Getwd()
		wr := bytes.NewBuffer(nil)
		tmpl.Execute(wr, map[string]string{
			"OS":   oa.OS,
			"Arch": oa.Arch,
			"Dir":  filepath.Base(cwd),
		})
		fmt.Printf("Building %s %s -> %s ...\n", oa.OS, oa.Arch, wr.String())
		//params := []string{"pack", "--rm", "--os", os.OS, "--arch", oa.Arch}
		cmd := exec.Command(os.Args[0], "pack",
			"-q", "--rm", "--os", oa.OS, "--arch", oa.Arch, "-o", wr.String())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
	outJson := ctx.String("json")
	if outJson == "" {
		return
	}
	comment, err := exec.Command("git", "log", "-1", "--oneline").Output()
	var commentStr string
	if err != nil {
		commentStr = os.Getenv("TRAVIS_COMMIT")
	} else {
		commentStr = string(comment)
	}
	vv := map[string]interface{}{
		"go_version":  runtime.Version(),
		"update_time": time.Now().Unix(),
		"format":      "zip",
		"comment":     commentStr,
		"builds":      osarches,
	}
	outfd, err := os.Create(outJson)
	if err != nil {
		log.Fatal(err)
	}
	defer outfd.Close()
	json.NewEncoder(outfd).Encode(vv)
}

func init() {
	cwd, _ := os.Getwd()
	program := filepath.Base(cwd)

	app.Name = "gopack"
	app.Usage = "Build and pack file into tgz or zip"

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: fmt.Sprintf("Generate %v file", RCFILE),
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Might rewrite config file",
				},
			},
			Action: actionInit,
		},
		{
			Name:  "all",
			Usage: fmt.Sprintf("Package all platform packages"),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "os",
					Usage: "Space-separated list of operating systems to build for",
					Value: DefaultPcfg.OS,
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Output path template",
					Value: "output/{{.Dir}}-{{.OS}}-{{.Arch}}.zip",
				},
				cli.StringFlag{
					Name:  "json",
					Usage: "Output json builds",
					Value: "",
				},
			},
			Action: actionAll,
		},
		{
			Name:  "pack",
			Usage: "Package file to zip or other format",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "os", EnvVar: "GOOS", Value: runtime.GOOS, Usage: "operation system"},
				cli.StringFlag{Name: "arch", EnvVar: "GOARCH", Value: runtime.GOARCH, Usage: "arch eg amd64|386|arm"},
				cli.StringFlag{Name: "output,o", Value: program + ".zip", Usage: "target file"},
				cli.StringFlag{Name: "gom", Value: "go", Usage: "go package manage program"},
				cli.BoolFlag{Name: "nobuild", Usage: "donot call go build when pack"},
				cli.BoolFlag{Name: "rm", Usage: "remove build files when done"},
				cli.BoolFlag{Name: "init", Usage: "generate sample .gobuild.yml"},
				cli.BoolFlag{Name: "quiet, q", Usage: "quiet console info"},
				cli.StringSliceFlag{Name: "add,a", Value: &cli.StringSlice{}, Usage: "add file"},
				//cli.StringFlag{Name: "depth", Value: "3", Usage: "depth of file to walk"},
			},
			Action: actionPack,
		},
		cmds.InstallFlag,
	}
	app.Version = VERSION
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "show debug information",
		},
	}
}

func main() {
	//app.Run(os.Args)
	app.RunAndExitOnError()
}
