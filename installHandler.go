package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/ssh"
)

func installHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	sshConfig := &ssh.ClientConfig{
		User: "student",
		Auth: []ssh.AuthMethod{ssh.Password("student")},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, _ := ssh.Dial("tcp", "157.179.132.136:22", sshConfig)
	session, _ := client.NewSession()
	var b bytes.Buffer
	session.Stdout = &b
	session.Run("ls -l")
	io.WriteString(w, b.String())
}
