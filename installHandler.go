package main

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func installHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	io.WriteString(w, "test")
}
