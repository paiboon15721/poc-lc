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

func getConnectionHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var err error
	var client *ssh.Client
	var session *ssh.Session
	sshConfig := &ssh.ClientConfig{
		User:    ps.ByName("username"),
		Auth:    []ssh.AuthMethod{ssh.Password(ps.ByName("password"))},
		Timeout: time.Duration(time.Millisecond * 2000),
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:22", ps.ByName("ip")), sshConfig)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Get server information
	session, _ = client.NewSession()
	var b bytes.Buffer
	session.Stdout = &b
	session.Run("lsb_release -a")
	io.WriteString(w, fmt.Sprintf("Connect success!\ncontent below is server information.\n-------------------------------------------------\n\n%s", b.String()))
}
