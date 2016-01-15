package example

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// NewServer makes a new example Server.
func NewServer() http.Handler {
	r := mux.NewRouter()
	r.Path("/hello").Methods("GET").HandlerFunc(handleHello)
	return r
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	fmt.Fprintf(w, "Hello %s.", q.Get("name"))
}
