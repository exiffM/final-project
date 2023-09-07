package executor

import (
	"bytes"
	"os/exec"
)

func RunCmd(command string) (string, error) {
	var outBuff, errBuff bytes.Buffer
	proc := exec.Command("bash", "-c", command)

	proc.Stdout = &outBuff
	proc.Stderr = &errBuff
	err := proc.Run()
	if err != nil {
		return "", err
	}
	return outBuff.String(), nil
}
