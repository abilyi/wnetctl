package sshclient

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"strings"
)

type SshClient interface {
	SetKey(keyPath string)
	Connect() error
	Execute(command string) error
	ExecuteInteractive(process InteractiveProcess) error
	Close() error
}

type InteractiveProcess interface {
	Command() []string
	Execute(stdin io.Writer, stdout io.Reader, stderr io.Reader) error
	//InSession(session *sshclient.Session) error
}

type sshClient struct {
	ip       string
	username string
	password string
	key      string
	client   *ssh.Client
}

type CommandsExecutionError struct {
	Index    int
	Command  string
	ExitCode int
	Stderr   string
}

func (this *CommandsExecutionError) Error() string {
	return fmt.Sprintf("Command \"%s\" execution failed with exit code %d", this.Command, this.ExitCode)
}

func NewCommandExecutionError(command string, exitCode int, stderr string) error {
	return &CommandsExecutionError{Command: command, ExitCode: exitCode, Stderr: stderr}
}

func (this *sshClient) SetKey(keyPath string) {
	this.key = keyPath
}

func (this *sshClient) Connect() error {
	var hostkey ssh.PublicKey
	config := &ssh.ClientConfig{
		User: this.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(this.password),
		},
		HostKeyCallback: ssh.FixedHostKey(hostkey),
	}
	client, err := ssh.Dial("tcp", this.ip+":22", config)
	if err != nil {
		return err
	}
	this.client = client
	return nil
}

func (this *sshClient) Execute(command string) error {
	session, err := this.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	return session.Run(command)
}

func (this *sshClient) ExecuteInteractive(process InteractiveProcess) error {
	session, err := this.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	pipes := make([]*os.File, 6)
	for i := 0; i < 3; i++ {
		pipes[i*2], pipes[i*2+1], err = os.Pipe()
		if err != nil {
			return err
		}
		defer pipes[i*2].Close()
		defer pipes[i*2+1].Close()
	}
	session.Stdin = pipes[0]
	session.Stdout = pipes[3]
	session.Stderr = pipes[5]
	err = session.Start(strings.Join(process.Command(), " "))
	if err != nil {
		return err
	}
	err = process.Execute(pipes[1], pipes[2], pipes[4])
	if err != nil {
		return err
	}
	return session.Wait()
}

func (this *sshClient) Close() error {
	return this.client.Close()
}

func NewSshClient(ip, login, password, sshKey string) SshClient {
	client := &sshClient{ip: ip, username: login, password: password, key: sshKey}
	return client
}
