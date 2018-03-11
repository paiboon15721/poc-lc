package main

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	router := httprouter.New()
	router.GET("/", indexHandler)
	router.POST("/", generateHandler)
	http.ListenAndServe(":3001", router)
}
