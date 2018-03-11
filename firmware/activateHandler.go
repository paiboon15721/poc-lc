package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	bolt "github.com/coreos/bbolt"
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

	var reqData reqBody
	json.NewDecoder(req.Body).Decode(&reqData)
	connectedHardwareID := strings.TrimSpace(reqData.HardwareID)
	if hardwareID == "" {
		http.Error(w, "hardwareID is required", 401)
		return
	}
	var exist bool
	db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte("ids")).Get([]byte(connectedHardwareID))
		if v != nil {
			exist = true
		}
		return nil
	})
	if !exist {
		var usage int
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("ids"))
			usage = b.Stats().KeyN
			return nil
		})
		if usage >= quotaTotal {
			http.Error(w, "run out of quota", 401)
			return
		}
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("ids"))
			if err := b.Put([]byte(connectedHardwareID), []byte("1")); err != nil {
				return fmt.Errorf("put: %s", err)
			}
			return nil
		})
	}
	json.NewEncoder(w).Encode(resBody{connectedHardwareID})
}
