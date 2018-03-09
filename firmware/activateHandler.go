package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func activateHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	type reqBody struct {
		HardwareID string `json:"hardwareID"`
		Nonce      string `json:"nonce"`
		Md         string `json:"md"`
	}

	type resBody struct {
		Md string `json:"md"`
	}
}
