package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

func getSSHClient(serverIP, serverUsername, serverPassword string) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User:    serverUsername,
		Auth:    []ssh.AuthMethod{ssh.Password(serverPassword)},
		Timeout: time.Duration(time.Millisecond * 2000),
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	return ssh.Dial("tcp", fmt.Sprintf("%s:22", serverIP), sshConfig)
}
