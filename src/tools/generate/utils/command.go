package utils

import (
	"bytes"
	"fmt"
	"github.com/hexya-erp/hexya/src/tools/logging"
	"os/exec"
)

var logger = logging.GetLogger("utils/command")

// RunCommand runs a shell command and logs the output.
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	logger.Info(fmt.Sprintf("Running command: %s %v", name, args))

	err := cmd.Run()
	if err != nil {
		logger.Error(fmt.Sprintf("Error running command: %s", stderr.String()))
		return fmt.Errorf("command execution failed: %v", err)
	}

	logger.Info(fmt.Sprintf("Command output: %s", out.String()))
	return nil
}
