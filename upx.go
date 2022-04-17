package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ergoapi/util/file"
	"github.com/sirupsen/logrus"
)

type (

	// Plugin defines the Docker plugin parameters.
	Plugin struct {
		Level   int // compress Level 1 fastest, 9 best compression
		Path    string
		Include string
		Exclude string
	}
)

// Exec executes the plugin step
func (p Plugin) Exec() error {
	cmd := commandInfo()
	if err := cmd.Run(); err != nil {
		logrus.Errorf("error running upx info command: %s", err)
		return err
	}
	if p.Path == "" {
		return fmt.Errorf("path is required")
	}
	if p.Level < 1 || p.Level > 9 {
		logrus.Warnf("level allow 1-9, auto change %d to 9", p.Level)
		p.Level = 9
	}
	var cmds []*exec.Cmd
	if file.IsDir(p.Path) {
		// dir compress
		files, err := file.DirFilesList(p.Path, p.Include, p.Exclude)
		if err != nil {
			return err
		}
		for _, file := range files {
			cmds = append(cmds, exec.Command("/usr/bin/upx", "-q", "-9", "-f", file))
		}
	} else {
		// file compress
		cmds = append(cmds, exec.Command("/usr/bin/upx", "-q", "-9", "-f", p.Path))
	}
	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// helper function to create the docker info command.
func commandInfo() *exec.Cmd {
	return exec.Command("/usr/bin/upx", "info")
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
