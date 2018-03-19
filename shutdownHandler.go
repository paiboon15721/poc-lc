package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func shutdownHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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

	// Shutdown server
	session, _ := client.NewSession()
	var b bytes.Buffer
	session.Stdout = &b
	session.Run(fmt.Sprintf("echo \"%s\" | sudo -S poweroff", serverPassword))
	io.WriteString(w, fmt.Sprintf("Shutdown server success!: %s", b.String()))
}
