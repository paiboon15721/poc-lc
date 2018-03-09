package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func infoHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	type quota struct {
		Total  int `json:"total"`
		Remain int `json:"remain"`
	}

	type info struct {
		Version    string `json:"version"`
		HardwareID string `json:"hardwareID"`
		Customer   string `json:"customer"`
		Quota      quota  `json:"quota"`
	}

	i := info{"0.0.1", hardwareID, customer, quota{quotaTotal, quotaTotal}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(i)
}
