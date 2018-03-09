package main

import (
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	jen "github.com/dave/jennifer/jen"
	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	router := httprouter.New()
	router.GET("/", index)
	router.POST("/", generate)
	http.ListenAndServe(":3001", router)
}

func generate(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if err := req.ParseForm(); err != nil {
		log.Fatalln(err)
	}
	f := jen.NewFile("main")
	customerName := req.PostFormValue("customerName")
	serverID := req.PostFormValue("serverID")
	totalLicense, _ := strconv.Atoi(req.PostFormValue("totalLicense"))
	f.Var().Add(jen.Id("customerName"), jen.Op("="), jen.Lit(customerName))
	f.Var().Add(jen.Id("serverID"), jen.Op("="), jen.Lit(serverID))
	f.Var().Add(jen.Id("totalLicense"), jen.Op("="), jen.Lit(totalLicense))
	f.Save("firmware/config.go")
	err := exec.Command("go", "build", "-o", "./firmware/firmware.exe", "./firmware").Run()
	if err != nil {
		log.Fatalln(err)
	}
	err = tpl.ExecuteTemplate(w, "gfw.html", "generate firmware success!")
	if err != nil {
		log.Fatalln(err)
	}
}

func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "gfw.html", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
