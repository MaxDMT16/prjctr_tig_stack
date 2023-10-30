package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"tig-stack/backend"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(
		w, `
          ##         .
    ## ## ##        ==
 ## ## ## ## ##    ===
/"""""""""""""""""\___/ ===
{                       /  ===-
\______ O           __/
 \    \         __/
  \____\_______/

	
Hello from Docker!

`,
	)
}


func main() {
	esService := backend.NewElasticsearch()
	
	h := backend.NewHandler(esService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handler)
	r.Get("/messages", h.SearchMessages)

	fmt.Println("Go backend started!")
	log.Fatal(http.ListenAndServe(":80", r))
}
