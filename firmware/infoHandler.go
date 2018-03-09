package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type quota struct {
	Total  int `json:"total"`
	Remain int `json:"remain"`
}

type information struct {
	Version    string `json:"version"`
	HardwareID string `json:"hardwareID"`
	Customer   string `json:"customer"`
	Quota      quota  `json:"quota"`
}

func info(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

}
