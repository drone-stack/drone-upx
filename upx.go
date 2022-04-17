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
		Debug   bool
	}
)

// Exec executes the plugin step
func (p Plugin) Exec() error {
	cmd := commandInfo()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
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
	p.Path = strings.TrimSuffix(p.Path, "/")
	level := fmt.Sprintf("-%d", p.Level)
	isdir := file.IsDir(p.Path)
	logrus.Debugf("path: %s, dir: %v", p.Path, isdir)
	if isdir {
		// dir compress
		files, err := file.DirFilesList(p.Path, p.Include, p.Exclude)
		if err != nil {
			return err
		}
		logrus.Debugf("path: %s, include: %s, exclude: %s, files: %v", p.Path, p.Include, p.Exclude, files)
		for _, f := range files {
			f = fmt.Sprintf("%s/%s", p.Path, f)
			if file.IsBinary(f) {
				cmds = append(cmds, exec.Command("/usr/bin/upx", "-q", level, "-f", f))
			}
		}
	} else {
		// file compress
		if file.IsBinary(p.Path) {
			cmds = append(cmds, exec.Command("/usr/bin/upx", "-q", level, "-f", p.Path))
		} else {
			return fmt.Errorf("%s is not a binary file", p.Path)
		}
	}
	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)
		if err := cmd.Run(); err != nil {
			return err
		}
		logrus.Infof("run [%s] compress success", strings.Join(cmd.Args, " "))
	}
	return nil
}

// helper function to create the docker info command.
func commandInfo() *exec.Cmd {
	return exec.Command("/usr/bin/upx", "-V")
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
