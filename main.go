package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
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
	OS   string
	Arch string
}

func actionAll(ctx *cli.Context) {
	ss := map[string][]string{
		"windows": {"amd64", "386"},
		"linux":   {"amd64", "386", "arm"},
		"darwin":  {"amd64"},
	}

	oses := ctx.StringSlice("os")
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
		cmd := exec.Command(os.Args[0], "pack",
			"-q", "--os", oa.OS, "--arch", oa.Arch, "-o", wr.String())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func init() {
	cwd, _ := os.Getwd()
	program := filepath.Base(cwd)

	app.Name = "gopack"
	app.Usage = "Build and pack file into tgz or zip"
	//app.Action = actionPack

	initOS := cli.StringSlice([]string{"linux", "darwin", "windows"})
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
				cli.StringSliceFlag{
					Name:  "os",
					Usage: "Operation system",
					Value: &initOS,
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Output path template",
					Value: "output/{{.Dir}}-{{.OS}}-{{.Arch}}.zip",
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
