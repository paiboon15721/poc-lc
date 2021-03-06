package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

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
	var session *ssh.Session
	client, err := getSSHClient(serverIP, serverUsername, serverPassword)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// Generate lcmgr.service file
	var tpl *template.Template
	var nf *os.File
	tpl, err = template.ParseFiles("lcmgr.service")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	nf, err = os.Create("firmware/lcmgr.service")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer nf.Close()
	err = tpl.Execute(nf, serverUsername)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Get hardwareID
	session, _ = client.NewSession()
	var b bytes.Buffer
	session.Stdout = &b
	session.Run("cat /proc/cpuinfo | grep Serial | cut -d ' ' -f 2")
	hardwareID := strings.TrimSuffix(b.String(), "\n")

	// Generate config.go file
	f := jen.NewFile("main")
	f.Var().Add(jen.Id("buildTime"), jen.Op("="), jen.Lit(time.Now().Format("02/01/2006 15:04:05")))
	f.Var().Add(jen.Id("customer"), jen.Op("="), jen.Lit(customer))
	f.Var().Add(jen.Id("hardwareID"), jen.Op("="), jen.Lit(hardwareID))
	f.Var().Add(jen.Id("quotaTotal"), jen.Op("="), jen.Lit(quotaTotal))
	f.Save("firmware/config.go")

	// Set go env for build firmware
	os.Setenv("GOARCH", "arm")
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
	session.Run(fmt.Sprintf("rm -rf lc && chmod +x firmware && mkdir -p lc && mv firmware lcmgr.service lc/ && echo \"%s\" | sudo -S cp lc/lcmgr.service /etc/systemd/system/ && echo \"%s\" | sudo -S systemctl enable lcmgr.service && echo \"%s\" | sudo -S systemctl restart lcmgr.service", serverPassword, serverPassword, serverPassword))
	io.WriteString(w, "Installed firmware success!")
}
