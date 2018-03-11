package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func indexHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "gfw.html", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
