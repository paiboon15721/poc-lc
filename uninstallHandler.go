package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func uninstallHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Get form body data
	// if err := req.ParseForm(); err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }
	// serverIP := req.PostFormValue("serverIP")
	// serverUsername := req.PostFormValue("serverUsername")
	// serverPassword := req.PostFormValue("serverPassword")
}
