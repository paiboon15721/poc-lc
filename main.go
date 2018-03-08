package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

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
	router.GET("/generate", testgenerate)
	router.POST("/generate", generate)
	http.ListenAndServe(":3001", router)
}

func generate(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	io.WriteString(w, "yo")
}

func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "gfw.html", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func testgenerate(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	f := jen.NewFile("main")
	f.Func().Id("main").Params().Block(
		jen.Qual("fmt", "Println").Call(jen.Lit("Hello, world")),
	)
	io.WriteString(w, fmt.Sprintf("%#v", f))
}
