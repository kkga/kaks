package kak

import (
	"fmt"
	"os"
	"syscall"
)

func Run(session string, kakArgs []string, fp *Filepath) error {
	kakExec, err := kakExec()
	if err != nil {
		return err
	}

	kakExecArgs := []string{kakExec}

	for _, a := range kakArgs {
		switch a {
		case "-c":
			kakExecArgs = append(kakExecArgs, "-c", session)
		default:
			return fmt.Errorf("unknown argument to Run: %s", a)
		}
	}

	if fp.Name != "" {
		kakExecArgs = append(kakExecArgs, fp.Name)

		if fp.Line != 0 {
			kakExecArgs = append(kakExecArgs, fmt.Sprintf("+%d:%d", fp.Line, fp.Column))
		}

	}

	execErr := syscall.Exec(kakExec, kakExecArgs, os.Environ())
	if execErr != nil {
		return execErr
	}

	return nil
}
