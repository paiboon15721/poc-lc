package main

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

func installHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var client *ssh.Client
	var session *ssh.Session
	sshConfig := &ssh.ClientConfig{
		User: "student",
		Auth: []ssh.AuthMethod{ssh.Password("student")},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, _ = ssh.Dial("tcp", "157.179.132.136:22", sshConfig)

	session, _ = client.NewSession()
	scp.CopyPath("./firmware/firmware", "firmware", session)

	session, _ = client.NewSession()
	scp.CopyPath("./firmware/lcmgr.service", "lcmgr.service", session)

	session, _ = client.NewSession()
	session.Run("chmod +x firmware && mkdir -p lc && mv firmware lcmgr.service lc/ && echo \"student\" | sudo -S cp lc/lcmgr.service /etc/systemd/system/ && echo \"student\" | sudo -S systemctl enable lcmgr.service && echo \"student\" | sudo -S systemctl start lcmgr.service")
	io.WriteString(w, "ok")
}
