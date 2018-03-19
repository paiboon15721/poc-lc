package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func uninstallHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Get form body data
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	serverIP := req.PostFormValue("serverIP")
	serverUsername := req.PostFormValue("serverUsername")
	serverPassword := req.PostFormValue("serverPassword")

	client, err := getSSHClient(serverIP, serverUsername, serverPassword)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Uninstall firmware
	session, _ := client.NewSession()
	var b bytes.Buffer
	session.Stdout = &b
	session.Run(fmt.Sprintf("rm -rf lc && echo \"%s\" | sudo -S systemctl stop lcmgr.service && echo \"%s\" | sudo -S systemctl disable lcmgr.service && echo \"%s\" | sudo -S rm -f /etc/systemd/system/lcmgr.service && echo \"%s\" | sudo -S systemctl daemon-reload && echo \"%s\" | sudo -S systemctl reset-failed", serverPassword, serverPassword, serverPassword, serverPassword, serverPassword))
	io.WriteString(w, fmt.Sprintf("Uninstall firmware success!: %s", b.String()))
}
