package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.ServeFiles("/app/*filepath", http.Dir("assets"))
	router.POST("/api/install", installHandler)
	router.GET("/api/scan-ip", scanIPHandler)
	http.ListenAndServe(":3001", router)
}
