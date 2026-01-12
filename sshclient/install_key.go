package sshclient

import (
	"io"
	"os"
)

type installSshKey string

func (this *installSshKey) Command() []string {
	return []string{"/usr/bin/tee", "-a", "/root/.sshclient/authorized_keys"}
}

func (this *installSshKey) Execute(stdin io.Writer, stdout io.Reader, stderr io.Reader) error {
	file, err := os.OpenFile(string(*this), os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	publicKey := make([]byte, fileInfo.Size())
	_, err = file.Read(publicKey)
	if err != nil {
		return err
	}
	_, err = stdin.Write(publicKey)
	return err
}

func NewInstallSshKey(publicKey string) InteractiveProcess {
	cmd := installSshKey(publicKey)
	return &cmd
}
