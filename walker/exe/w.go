package exe

import (
	"os"
	"os/exec"
	"strings"
)

func NewExe(wDir string, command string) *exec.Cmd {
	var cs = strings.Split(command, " ")
	var cmd = exec.Command(cs[0], cs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = wDir
	return cmd
}

func RunExe(wDir string, command string) error {
	var cmd = NewExe(wDir, command)
	return cmd.Run()
}
