package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	bolt "github.com/coreos/bbolt"
	"github.com/julienschmidt/httprouter"
)

var db *bolt.DB

func main() {
	// Check valid hardwareID
	var err error
	cmd := exec.Command("sudo", "cat", "/sys/class/dmi/id/product_uuid")
	var b bytes.Buffer
	cmd.Stdout = &b
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	realHardwareID := strings.TrimSuffix(b.String(), "\n")
	if realHardwareID != hardwareID {
		panic("hardwareID invalid!")
	}

	// Initial bolt database
	db, err = bolt.Open("lc.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("ids"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	// Setup router
	router := httprouter.New()
	router.GET("/info", infoHandler)
	router.POST("/activate", activateHandler)
	http.ListenAndServe(":3001", router)
}
