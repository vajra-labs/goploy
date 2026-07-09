package shell

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

type ExecError struct {
	Message  string
	Command  string
	Stdout   string
	Stderr   string
	ExitCode int
	ServerID *string
	Err      error
}

func (e *ExecError) Error() string {
	return e.Message
}

// IsRemote returns true if this error is from a remote execution.
func (e *ExecError) IsRemote() bool {
	return e.ServerID != nil
}

// Get a formatted error message with all details
func (e *ExecError) DetailedMessage() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Command: %s", e.Command))
	parts = append(parts, fmt.Sprintf("Exit Code: %d", e.ExitCode))
	if e.ServerID != nil {
		parts = append(parts, fmt.Sprintf("Server ID: %s", *e.ServerID))
	} else {
		parts = append(parts, "Location: Local")
	}
	if e.Stderr != "" {
		parts = append(parts, "Stderr: "+e.Stderr)
	}
	if e.Stdout != "" {
		parts = append(parts, "Stdout: "+e.Stdout)
	}
	return e.Error() + "\n" + strings.Join(parts, "\n")
}

func newExecError(command, stdout, stderr string, err error) *ExecError {
	execErr := &ExecError{
		Message: fmt.Sprintf("command execution failed: %v", err),
		Command: command,
		Stdout:  stdout,
		Stderr:  stderr,
		Err:     err,
	}
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		execErr.ExitCode = exitErr.ExitCode()
	}
	return execErr
}

func newSSHExecError(command, stdout, stderr string, err error, serverId string) *ExecError {
	execErr := &ExecError{
		Message:  fmt.Sprintf("remote command execution failed: %v", err),
		Command:  command,
		Stdout:   stdout,
		Stderr:   stderr,
		ServerID: &serverId,
		Err:      err,
	}
	// golang.org/x/crypto/ssh exposes exit status via *ssh.ExitError
	if sshExitErr, ok := errors.AsType[*ssh.ExitError](err); ok {
		execErr.ExitCode = sshExitErr.ExitStatus()
	}
	return execErr
}
