package sshclient

import (
	"bufio"
	"errors"
	"io"
)

func NewPasswd(user, password, newPassword string) InteractiveProcess {
	return &passwd{user, password, newPassword}
}

type passwd struct {
	user        string
	oldPassword string
	newPassword string
}

func (this *passwd) Command() []string {
	return []string{"/sbin/passwd", this.user}
}

func (this *passwd) Execute(stdin io.Writer, stdout io.Reader, stderr io.Reader) error {
	out := bufio.NewReader(stdout)

	prompt, err := out.ReadString('\n')
	if err != nil {
		return err
	}
	if prompt != "Changing password for "+this.user {
		return errors.New("Unexpected prompt from passwd")
	}
	prompt, err = out.ReadString(':')
	if err != nil {
		return err
	}
	if prompt != "New password:" {
		return errors.New("Unexpected prompt from passwd")
	}
	_, err = stdin.Write([]byte(this.newPassword))
	return err
}
