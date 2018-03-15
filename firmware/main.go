package main

import (
	"fmt"
	"net/http"

	bolt "github.com/coreos/bbolt"
	"github.com/julienschmidt/httprouter"
)

var db *bolt.DB

func main() {
	// Check valid hardwareID
	var err error
	// realHardwareIDbyte, err := ioutil.ReadFile("/sys/class/dmi/id/product_uuid")
	// realHardwareID := strings.TrimSuffix(string(realHardwareIDbyte), "\n")
	// if err != nil {
	// 	panic(err)
	// }
	// if realHardwareID != hardwareID {
	// 	panic("hardwareID invalid!")
	// }

	// Initial bolt database
	db, err = bolt.Open("lc.db", 0644, nil)
	if err != nil {
		panic(err)
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
