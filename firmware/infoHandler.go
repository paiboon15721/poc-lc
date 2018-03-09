package main

import (
	"encoding/json"
	"net/http"

	bolt "github.com/coreos/bbolt"
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
	var usage int
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ids"))
		usage = b.Stats().KeyN
		return nil
	})
	i := info{"0.0.1", hardwareID, customer, quota{quotaTotal, quotaTotal - usage}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(i)
}
