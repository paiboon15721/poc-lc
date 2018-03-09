package main

import (
	"fmt"
	"log"
	"net/http"

	bolt "github.com/coreos/bbolt"
	"github.com/julienschmidt/httprouter"
)

var db *bolt.DB

func main() {
	var err error
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
	router := httprouter.New()
	router.GET("/info", infoHandler)
	router.POST("/activate", activateHandler)
	http.ListenAndServe(":3001", router)
}
