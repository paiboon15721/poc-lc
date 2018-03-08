package main

import (
	"fmt"
	"io"
	"net/http"

	jen "github.com/dave/jennifer/jen"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/", generate)
	http.ListenAndServe(":3001", router)
}

func generate(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	f := jen.NewFile("main")
	f.Func().Id("main").Params().Block(
		jen.Qual("fmt", "Println").Call(jen.Lit("Hello, world")),
	)
	io.WriteString(w, fmt.Sprintf("%#v", f))
}
