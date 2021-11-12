package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type CmdOpt func(cmd *exec.Cmd)

func WithEnv(env []string) CmdOpt {
	return func(cmd *exec.Cmd) {
		cmd.Env = append(cmd.Env, env...)
	}
}

func WithDir(dir string) CmdOpt {
	return func(cmd *exec.Cmd) {
		cmd.Dir = dir
	}
}

func Run(command string, args []string, stdin io.Reader, stdout, stderr io.Writer, extraOpts ...CmdOpt) (exitCode int, err error) {
	cmd := exec.Command(command, args...)

	cmd.Env = os.Environ()
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	if stderr == nil {
		stderr = new(strings.Builder)
	}
	cmd.Stderr = stderr

	for _, fn := range extraOpts {
		fn(cmd)
	}

	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return 1, err
		}
		if exitErr := new(exec.ExitError); errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
		if b, ok := stderr.(*strings.Builder); ok {
			return exitCode, fmt.Errorf(b.String())
		}
		return 1, err
	}
	return 0, nil
}
