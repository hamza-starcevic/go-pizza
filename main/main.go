// ! This file contains the code for the creation of the server,
// ! the registration of the routes and their respective handlers
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const port = ":8070"

func main() {

	//* A router from the gorilla/mux package is created
	r := mux.NewRouter()
	//* The router is used to register the routes and their respective handlers
	//* A middleware to handle the CORS policy is also registered
	r.Use(mux.CORSMethodMiddleware(r))

	//* A post method for pizza creation is registered
	r.HandleFunc("/pizza", createPizza).Methods("POST")

	//* A delete method for pizza deletion is registered
	r.HandleFunc("/pizza/{id}", deletePizzaById).Methods("DELETE")

	//* A get method for all pizza retrieval is registered
	r.HandleFunc("/pizza", fetchPizza).Methods("GET")

	//* A get method for pizza retrieval by id is registered
	r.HandleFunc("/pizza/{id}", fetchPizzaById).Methods("GET")

	//* A get method for pizza retrieval by category is registered
	r.HandleFunc("/pizzaByCategory/{category}", fetchPizzaByCategory).Methods("GET")

	//* A get method for pizza retrieval by name is registered
	r.HandleFunc("/searchPizza/{name}", searchPizza).Methods("GET")

	//* A put method for pizza rating update is registered
	r.HandleFunc("/pizza/{id}/{column}/{value}", updatePizzaRatingById).Methods("PUT")

	//* A put method for pizza update is registered
	r.HandleFunc("/pizza", updatePizzaById).Methods("PUT")

	//* A get method for serving the html file is registered
	r.HandleFunc("/", serveHTMLfunc)

	//* The server is started on the specified port
	log.Printf("Server started on port %s", port)
	http.ListenAndServe(port, r)
}

// * Short handler function for serving the html to a minimal pizza adding site
func serveHTMLfunc(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
