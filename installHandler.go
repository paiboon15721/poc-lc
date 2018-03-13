package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	jen "github.com/dave/jennifer/jen"
	"github.com/julienschmidt/httprouter"
	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

func installHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Get form body data
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	serverIP := req.PostFormValue("serverIP")
	serverUsername := req.PostFormValue("serverUsername")
	serverPassword := req.PostFormValue("serverPassword")
	customer := req.PostFormValue("customer")
	quotaTotal, _ := strconv.Atoi(req.PostFormValue("quotaTotal"))

	// Init ssh session
	var err error
	var client *ssh.Client
	var session *ssh.Session
	sshConfig := &ssh.ClientConfig{
		User: serverUsername,
		Auth: []ssh.AuthMethod{ssh.Password(serverPassword)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:22", serverIP), sshConfig)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// Get hardwareID
	session, _ = client.NewSession()
	var b bytes.Buffer
	session.Stdout = &b
	session.Run(fmt.Sprintf("echo \"%s\" | sudo -S cat /sys/class/dmi/id/product_uuid", serverPassword))
	hardwareID := strings.TrimSuffix(b.String(), "\n")

	// Generate config.go file
	f := jen.NewFile("main")
	f.Var().Add(jen.Id("customer"), jen.Op("="), jen.Lit(customer))
	f.Var().Add(jen.Id("hardwareID"), jen.Op("="), jen.Lit(hardwareID))
	f.Var().Add(jen.Id("quotaTotal"), jen.Op("="), jen.Lit(quotaTotal))
	f.Save("firmware/config.go")

	// Set go env for build firmware
	os.Setenv("GOARCH", "386")
	os.Setenv("GOOS", "linux")

	// Build firmware
	err = exec.Command("go", "build", "-o", "./firmware/firmware", "./firmware").Run()
	if err != nil {
		http.Error(w, "build firmware fail!", 500)
		return
	}

	// Restore go env for develop
	os.Setenv("GOARCH", "amd64")
	os.Setenv("GOOS", "windows")

	// Scp firmware to server
	session, _ = client.NewSession()
	scp.CopyPath("./firmware/firmware", "firmware", session)
	session, _ = client.NewSession()
	scp.CopyPath("./firmware/lcmgr.service", "lcmgr.service", session)

	// Start license manager service
	session, _ = client.NewSession()
	session.Run("chmod +x firmware && mkdir -p lc && mv firmware lcmgr.service lc/ && echo \"student\" | sudo -S cp lc/lcmgr.service /etc/systemd/system/ && echo \"student\" | sudo -S systemctl enable lcmgr.service && echo \"student\" | sudo -S systemctl restart lcmgr.service")
	io.WriteString(w, "Installed firmware success!")
}
