package main

import (
	"net/http"
	"os"

	"github.com/Fontinalis/gogcs"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	gcs, err := gogcs.New()
	if err != nil {
		panic(err)
	}
	r.PathPrefix("/").Handler(gcs)
	http.ListenAndServe(":3698", handlers.CombinedLoggingHandler(os.Stdout, r))
}
