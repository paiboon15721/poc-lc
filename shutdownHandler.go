package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/ssh"
)

func shutdownHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var err error
	serverIP := req.PostFormValue("serverIP")
	serverUsername := req.PostFormValue("serverUsername")
	serverPassword := req.PostFormValue("serverPassword")
	var client *ssh.Client
	var session *ssh.Session
	sshConfig := &ssh.ClientConfig{
		User:    serverUsername,
		Auth:    []ssh.AuthMethod{ssh.Password(serverPassword)},
		Timeout: time.Duration(time.Millisecond * 2000),
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:22", serverIP), sshConfig)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Shutdown server
	session, _ = client.NewSession()
	var b bytes.Buffer
	session.Stdout = &b
	session.Run(fmt.Sprintf("echo \"%s\" | sudo -S poweroff", serverPassword))
	io.WriteString(w, fmt.Sprintf("Shutdown server success!: %s", b.String()))
}
