package runner

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func Run(command string) error {
	fmt.Fprintf(os.Stderr, "$ %s\n", command)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
