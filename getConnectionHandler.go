package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func getConnectionHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	client, err := getSSHClient(
		ps.ByName("serverIP"),
		ps.ByName("serverUsername"),
		ps.ByName("serverPassword"),
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Get server information
	session, _ := client.NewSession()
	var b bytes.Buffer
	session.Stdout = &b
	session.Run("lsb_release -a")
	io.WriteString(w, fmt.Sprintf("Connect success!\nThe contents below are server information.\n-------------------------------------------------\n\n%s", b.String()))
}
