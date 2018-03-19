package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func getInfoHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	client := http.Client{
		Timeout: time.Second,
	}
	resp, err := client.Get(fmt.Sprintf("http://%s:3001/info", ps.ByName("serverIP")))
	if err != nil {
		http.Error(w, "Service unavailable", 503)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}
