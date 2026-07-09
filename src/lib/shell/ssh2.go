package shell

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// Global connection pool — keyed by serverId.
var pool sync.Map // map[string]*ssh.Client

// SSHConfig holds the connection parameters for a remote server.
type SSHConfig struct {
	Host       string
	Port       string
	User       string
	PrivateKey string
}

// SSHClient is a lightweight handle for a registered server.
// Obtain one via GetSSHClient after registering with NewSSH.
type SSHClient struct {
	serverId string
	conn     *ssh.Client
}

// NewSSH dials the remote server, establishes an SSH connection, and stores
// it in the global pool under serverId. Calling NewSSH again with the same
// serverId replaces the existing connection.
func NewSSH(serverId string, cfg SSHConfig) error {
	signer, err := ssh.ParsePrivateKey([]byte(cfg.PrivateKey))
	if err != nil {
		return &ExecError{
			Message: fmt.Sprintf("invalid SSH private key: %v", err),
			Err:     err,
		}
	}
	client, err := ssh.Dial("tcp", cfg.Host+":"+cfg.Port, &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	})
	if err != nil {
		msg := err.Error()
		isAuthFailure := strings.Contains(msg, "unable to authenticate") ||
			strings.Contains(msg, "no supported methods remain") ||
			strings.Contains(msg, "ssh: handshake failed")
		if isAuthFailure {
			technicalDetail := fmt.Sprintf("Error: %v", err)
			friendlyMessage := strings.Join([]string{
				"",
				"❌ Couldn't connect to your server — the SSH key was not accepted.",
				"",
				"This usually means the key doesn't match what's on the server, or the key format is invalid.",
				"",
				"Technical details: " + technicalDetail,
				"",
				"💡 Hints:",
				"  • Check that the SSH key you added is the same one installed on the server (e.g. in ~/.ssh/authorized_keys).",
				"  • Try generating a new SSH key and add only the public key to the server, then try again.",
				"  • Make sure the correct user and port are configured for the server.",
			}, "\n")
			return &ExecError{
				Message:  "Authentication failed: Invalid SSH private key.\n" + friendlyMessage,
				ServerID: &serverId,
				Err:      err,
			}
		}
		return &ExecError{
			Message:  fmt.Sprintf("SSH connection error: %v", err),
			ServerID: &serverId,
			Err:      err,
		}
	}
	pool.Store(serverId, client)
	return nil
}

// GetSSHClient retrieves a registered SSHClient handle from the pool.
// Returns an error if NewSSH has not been called for this serverId.
func GetSSHClient(serverId string) (*SSHClient, error) {
	val, ok := pool.Load(serverId)
	if !ok {
		return nil, fmt.Errorf("server %q not connected", serverId)
	}
	return &SSHClient{
		serverId: serverId,
		conn:     val.(*ssh.Client),
	}, nil
}

// Exec runs a command on the remote server and returns a channel that receives
// the result once execution completes.
//
// Each call opens a new SSH session on the existing TCP connection — lightweight
// and safe to call concurrently. If onData is set, output chunks are forwarded
// in real-time. The context is honoured — cancelling it closes the session.
func (c *SSHClient) Exec(ctx context.Context, command string, onData func(string)) <-chan ExecResult {
	ch := make(chan ExecResult, 1)
	go func() {
		defer close(ch)

		session, err := c.conn.NewSession()
		if err != nil {
			ch <- ExecResult{Err: fmt.Errorf("ssh new session: %w", err)}
			return
		}
		defer session.Close()

		done := make(chan struct{})
		defer close(done)
		go func() {
			select {
			case <-ctx.Done():
				session.Close()
			case <-done:
			}
		}()

		if onData != nil {
			ch <- sshExecStream(session, c.serverId, command, onData)
		} else {
			ch <- sshExecSimple(session, c.serverId, command)
		}
	}()
	return ch
}

// Close removes the client from the pool and closes the underlying connection.
func (c *SSHClient) Close() error {
	pool.Delete(c.serverId)
	return c.conn.Close()
}

// SSHCloseAll closes every connection in the pool.
func SSHCloseAll() {
	pool.Range(func(key, val any) bool {
		_ = val.(*ssh.Client).Close()
		pool.Delete(key)
		return true
	})
}

// sshExecSimple runs the command and captures stdout/stderr into buffers.
// Used when no streaming callback is provided.
func sshExecSimple(session *ssh.Session, serverId, command string) ExecResult {
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(command); err != nil {
		return ExecResult{Err: newSSHExecError(
			command, stdout.String(),
			stderr.String(), err,
			serverId,
		)}
	}
	return ExecResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
}

// sshExecStream runs the command and forwards output to onData if set.
func sshExecStream(session *ssh.Session, serverId, command string, onData func(string)) ExecResult {
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return ExecResult{Err: fmt.Errorf("ssh stdout pipe: %w", err)}
	}
	stderrPipe, err := session.StderrPipe()
	if err != nil {
		return ExecResult{Err: fmt.Errorf("ssh stderr pipe: %w", err)}
	}

	var stdout, stderr bytes.Buffer
	stdoutWriter := &streamWriter{buf: &stdout, onData: onData}
	stderrWriter := &streamWriter{buf: &stderr, onData: onData}

	if err := session.Start(command); err != nil {
		return ExecResult{Err: newSSHExecError(
			command, "", "",
			err, serverId,
		)}
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); io.Copy(stdoutWriter, stdoutPipe) }()
	go func() { defer wg.Done(); io.Copy(stderrWriter, stderrPipe) }()

	err = session.Wait()
	wg.Wait()

	if err != nil {
		return ExecResult{Err: newSSHExecError(
			command, stdout.String(),
			stderr.String(), err,
			serverId,
		)}
	}
	return ExecResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
}
