// ! This file contains the code for the creation of the server,
// ! the registration of the routes and their respective handlers
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	auth "github.com/hamza-starcevic/pizzaapi/pkg/authPizza"
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

	//! A delete method for pizza deletion is registered
	r.HandleFunc("/pizza/{id}", deletePizzaById).Methods("DELETE")

	//* A get method for all pizza retrieval is registered
	r.HandleFunc("/pizza", fetchPizza).Methods("GET")

	//* A get method for pizza retrieval by id is registered
	r.HandleFunc("/pizza/{id}", fetchPizzaById).Methods("GET")

	//* A get method for pizza retrieval by category is registered
	r.HandleFunc("/pizzaByCategory/{category}", fetchPizzaByCategory).Methods("GET")

	//* A get method for pizza retrieval by name is registered
	r.HandleFunc("/searchPizza/{name}", searchPizza).Methods("GET")

	//? A put method for pizza rating update is registered
	r.HandleFunc("/pizza/{id}/{column}/{value}", updatePizzaRatingById).Methods("PUT")

	//? A put method for pizza update is registered
	r.HandleFunc("/pizza/{id}", updatePizzaById).Methods("PUT")

	//* A get method for serving the html file is registered
	r.HandleFunc("/", serveHTMLfunc)

	//* A post method for user registration is registered
	r.HandleFunc("/auth/register", auth.RegisterUser).Methods("POST")

	//* A post method for user login is registered
	r.HandleFunc("/auth/login", auth.LoginUser).Methods("POST")

	//* A post method for order creation is registered
	r.HandleFunc("/order", createOrder).Methods("POST")

	//* A get method for order retrieval by id is registered
	r.HandleFunc("/order", getOrders).Methods("GET")

	//* A get method for order retrieval by id is registered
	r.HandleFunc("/orderhistory", getOrdersbyId).Methods("GET")

	//! A delete method for order deletion is registered
	r.HandleFunc("/order/{id}", deleteOrderById).Methods("DELETE")

	//* A get method for fetching all users is registered
	r.HandleFunc("/users", fetchUsers).Methods("GET")

	//! A delete method for deleting a user is registered
	r.HandleFunc("/users/{id}", deleteUserById).Methods("DELETE")

	//? A put method for updating a user is registered
	r.HandleFunc("/users/{id}", updateUserById).Methods("PUT")

	//* A route for serving the html file is registered
	r.HandleFunc("/test/{token}", verifyTokenTest).Methods("GET")

	//* The server is started on the specified port
	log.Printf("Server started on port %s", port)
	logMiddleware := handlers.LoggingHandler(os.Stdout, r)
	http.ListenAndServe(port, logMiddleware)
}

// * Short handler function for serving the html to a minimal pizza adding site
func serveHTMLfunc(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
