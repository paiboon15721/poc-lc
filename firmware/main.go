package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"unicode/utf8"

	bolt "github.com/coreos/bbolt"
	"github.com/julienschmidt/httprouter"
)

var db *bolt.DB

func scanHardwareID(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if r == 'S' {
			if string(data[start:start+6]) == "Serial" {
				start += 6
				break
			}
		}
	}
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if r == '\n' {
			return i + width, data[start:i], nil
		}
	}
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	return start, nil, nil
}

func main() {
	// Check valid hardwareID
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(scanHardwareID)
	scanner.Scan()
	realHardwareID := scanner.Text()
	realHardwareID = realHardwareID[len(realHardwareID)-16:]
	if realHardwareID != hardwareID {
		panic("hardwareID invalid!")
	}

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
