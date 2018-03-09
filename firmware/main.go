package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/info", info)
	http.ListenAndServe(":3000", router)
}
