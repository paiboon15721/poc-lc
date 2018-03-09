package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/info", infoHandler)
	router.POST("/activate", activateHandler)
	http.ListenAndServe(":3000", router)
}
