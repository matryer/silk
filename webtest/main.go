package main

import (
	"log"
	"net/http"

	"github.com/matryer/silk/testutil"
)

func main() {

	http.Handle("/", testutil.EchoHandler())

	err := http.ListenAndServe(":9080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
