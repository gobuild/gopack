package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/qiniu/log"
)

var InstallFlag = cli.Command{
	Name:  "install",
	Usage: fmt.Sprintf("Install new binary"),
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug log",
		},
	},
	Action: InstallAction,
}

var GOBIN string

// Determine binary path
func init() {
	var paths []string
	gobin := os.Getenv("GOBIN")
	if gobin != "" {
		paths = filepath.SplitList(gobin)
	}
	if len(paths) > 0 {
		GOBIN = paths[0]
		return
	}
	paths = filepath.SplitList(os.Getenv("GOPATH"))
	if len(paths) > 0 {
		GOBIN = filepath.Join(paths[0], "bin")
		return
	}
	log.Fatal("Make sure you set env GOPATH or GOBIN")

	os.ExpandEnv("$HOME/.gopack/src")
}

func getInsPath(names ...string) string {
	base := os.ExpandEnv("$HOME/.gopack/")
	names = append([]string{base}, names...)
	target := filepath.Join(names...)
	os.MkdirAll(filepath.Dir(target), 0755)
	return target
}

// http://dn-gobuild5.qbox.me/gorelease/gorelease/gorelease/master/darwin-amd64/gorelease.zip
// Need to update
func downloadSource(name string) (dest string, err error) {
	osarch := runtime.GOOS + "-" + runtime.GOARCH
	url := fmt.Sprintf("http://dn-gobuild5.qbox.me/gorelease/codeskyblue/%s/master/%s/%s.zip",
		name, osarch, name)
	log.Debug("download:", url)
	dest = getInsPath("src", fmt.Sprintf("%s.zip", name))
	cmd := exec.Command("curl", "-fSL", url, "-o", dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}

func deployPackage(pkgName, path string, binDir string) error {
	cmd := exec.Command("unzip", "-o", "-d", getInsPath("opt", pkgName), path)
	log.Debug("zip command:", cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	// for linux and darwin
	os.Remove(filepath.Join(GOBIN, pkgName))
	return os.Symlink(filepath.Join("../opt/", pkgName, pkgName), filepath.Join(GOBIN, pkgName)) //getInsPath("bin", pkgName))
}

func InstallAction(c *cli.Context) {
	if c.Bool("debug") {
		log.SetOutputLevel(log.Ldebug)
	}
	if len(c.Args()) < 1 {
		log.Fatal("Need at lease one argument")
	}
	log.Println(GOBIN)
	pkgName := c.Args().First()
	dest, err := downloadSource(pkgName)
	if err != nil {
		log.Fatal(err)
	}
	err = deployPackage(pkgName, dest, GOBIN)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Program [%s] installed\n", pkgName)
}
