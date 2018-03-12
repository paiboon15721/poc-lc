package main

import (
	"log"
	"net/http"
	"os/exec"
	"strconv"

	jen "github.com/dave/jennifer/jen"
	"github.com/julienschmidt/httprouter"
)

func generateHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if err := req.ParseForm(); err != nil {
		log.Fatalln(err)
	}
	f := jen.NewFile("main")
	customer := req.PostFormValue("customer")
	hardwareID := req.PostFormValue("hardwareID")
	quotaTotal, _ := strconv.Atoi(req.PostFormValue("quotaTotal"))
	f.Var().Add(jen.Id("customer"), jen.Op("="), jen.Lit(customer))
	f.Var().Add(jen.Id("hardwareID"), jen.Op("="), jen.Lit(hardwareID))
	f.Var().Add(jen.Id("quotaTotal"), jen.Op("="), jen.Lit(quotaTotal))
	f.Save("firmware/config.go")
	err := exec.Command("go", "build", "-o", "./firmware/firmware.exe", "./firmware").Run()
	if err != nil {
		log.Fatalln(err)
	}
}
