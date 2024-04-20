package sshclient

import (
	"fmt"
	"kubehostwarden/types"

	"golang.org/x/crypto/ssh"
)

func NewSSHClient(sshInfo types.SSHInfo) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: sshInfo.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshInfo.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshInfo.EndPoint, sshInfo.Port), sshConfig)
	if err != nil {
		return nil, err
	}

	return sshClient, nil
}
