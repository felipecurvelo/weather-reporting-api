package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func FirstEndpoint(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome! This is the first endpoint working!\n")
}

func main() {
	router := httprouter.New()
	router.GET("/", FirstEndpoint)
	log.Fatal(http.ListenAndServe(":8080", router))
}
