package cmds

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/codegangsta/cli"
)

var BuildFlag = cli.Command{
	Name:  "build",
	Usage: "Build with version code",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "version",
			Usage: "set version name",
		},
	},
	Action: BuildAction,
}

func BuildAction(c *cli.Context) {
	ver := c.String("version")
	if ver == "" {
		ver = time.Now().Format("2006-01-02_15:04:05")
	}
	cmd := exec.Command("go", "build", "-ldflags",
		fmt.Sprintf("-X main.VERSION=%s", strconv.Quote(ver)))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
