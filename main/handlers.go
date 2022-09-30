// ! This file contains the functions(handlers) for the user routes
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// * Struct used for indicating the status of an operation
type successJson struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

// * fetchPizza returns an array with all pizzas from the database
func fetchPizza(w http.ResponseWriter, r *http.Request) {

	//* establish connection to the database
	db := Connect()
	defer db.Close()

	//* initialize an array of pizzas(structs of type Pizza)
	var pizzaArr []Pizza

	//* query the database for all pizzas
	result, err := db.Query("SELECT * FROM Pizza")

	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	//* iterate through the result and append each pizza to the array
	for result.Next() {
		var pizza Pizza
		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			panic(err.Error())
		}
		pizzaArr = append(pizzaArr, pizza)
	}

	//* set the content type to json and encode the array of pizzas
	log.Println("All pizzas fetched")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizzaArr)
}

// * fetchPizzaByID returns a pizza with the given ID from the database
func fetchPizzaById(w http.ResponseWriter, r *http.Request) {

	//* establish connection to the database
	db := Connect()
	defer db.Close()

	//* mux.Vars(r) returns a map of the variables in the url
	mux.Vars(r)

	//* initialize a pizza(struct of type Pizza)
	var pizza Pizza

	//* query the database for the pizza with the given ID, where the ID is the first variable in the url
	result, err := db.Query("SELECT * FROM Pizza WHERE ID = ?", mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	//* iterate through the result and add the values to the pizza variable
	for result.Next() {
		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			panic(err.Error())
		}
	}

	//* set the content type to json and encode the pizza
	log.Println("Pizza searched: ", mux.Vars(r)["id"], pizza.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizza)
}

// * createPizza creates a new entry in the database
func createPizza(w http.ResponseWriter, r *http.Request) {

	//* establish connection to the database
	db := Connect()
	defer db.Close()

	//* initialize a pizza(struct of type Pizza), then decode the request body into the pizza variable
	var pizza Pizza
	json.NewDecoder(r.Body).Decode(&pizza)

	//* generate the DateAdded value
	pizza.DateAdded = time.Now().Format("2006-01-02 15:04:05")

	//* generate the ID value, which is a random UUID(a string of 36 characters which is unique and random)
	ID := uuid.New().String()
	pizza.ID = ID

	//* query the database to insert the pizza into the database using the values from the pizza variable
	_, err := db.Exec(`INSERT INTO
							Pizza (ID, Picture, Name, Price, Description, Category, Ingredients, Rating, Date_added)
					   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		pizza.ID, pizza.Picture, pizza.Name, pizza.Price, pizza.Description, pizza.Category, pizza.Ingredients, pizza.Rating, pizza.DateAdded)

	if err != nil {
		panic(err.Error())
	}

	//* set the content type to json and encode the pizza struct that was just created
	log.Println("Pizza created: ", pizza.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizza)
}

// * fetchPizzaByCategory returns an array with all the pizzas in the given category
func fetchPizzaByCategory(w http.ResponseWriter, r *http.Request) {

	//* establish connection to the database
	db := Connect()
	defer db.Close()

	//* mux.Vars(r) returns a map of the variables in the url
	mux.Vars(r)

	//* initialize an array of pizzas(structs of type Pizza)
	var pizzaArr []Pizza

	//* query the database for all pizzas in the given category, where the category is the first variable in the url
	result, err := db.Query("SELECT * FROM Pizza WHERE Category = ?", mux.Vars(r)["category"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	//* iterate through the result and append each pizza to the array
	for result.Next() {
		var pizza Pizza

		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			panic(err.Error())
		}
		pizzaArr = append(pizzaArr, pizza)
	}

	//* set the content type to json and encode the array of pizzas
	log.Println("Pizzas fetched by category: ", mux.Vars(r)["category"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizzaArr)
}

// * searchPizza returns an array with all the pizzas that contain the given string in their name
func searchPizza(w http.ResponseWriter, r *http.Request) {
	//* establish connection to the database
	db := Connect()
	defer db.Close()

	//* mux.Vars(r) returns a map of the variables in the url, then retrieves the search string from the map
	mux.Vars(r)
	name := mux.Vars(r)["name"]

	//* initialize an array of pizzas(structs of type Pizza)
	var pizzaArr []Pizza

	//* query the database for all pizzas that contain the given string in their name
	result, err := db.Query("SELECT * FROM Pizza WHERE Name LIKE ?", "%"+mux.Vars(r)["name"]+"%")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	//* iterate through the result and append each pizza to the array
	for result.Next() {
		var pizza Pizza

		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			panic(err.Error())
		}
		pizzaArr = append(pizzaArr, pizza)
	}

	//* set the content type to json and encode the array of pizzas
	log.Println("Pizzas searched by term: ", name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizzaArr)
}

// * deletePizzaById deletes a pizza with the given ID from the database
func deletePizzaById(w http.ResponseWriter, r *http.Request) {

	//* establish connection to the database
	db := Connect()
	defer db.Close()

	//* mux.Vars(r) returns a map of the variables in the url, then retrieves the ID from the map
	mux.Vars(r)

	//* query the database to delete the pizza with the given ID
	result, err := db.Exec("DELETE FROM Pizza WHERE ID = ?", mux.Vars(r)["id"])

	if err != nil {
		panic(err.Error())
	}

	//* set the content type to json and encode successJson struct depending on whether the pizza was deleted or not
	effect, _ := result.RowsAffected()
	w.Header().Set("Content-Type", "application/json")

	if effect == 0 {
		w.WriteHeader(http.StatusNotFound)
		var fail = successJson{Success: false, ID: mux.Vars(r)["id"]}
		json.NewEncoder(w).Encode(fail)

	} else {
		w.WriteHeader(http.StatusOK)
		var success = successJson{Success: true, ID: mux.Vars(r)["id"]}
		json.NewEncoder(w).Encode(success)
	}

	log.Println("Pizza Deleted: ", mux.Vars(r)["id"])
}

// * updatePizzaRatingById updates the rating of a pizza with the given ID from the database
func updatePizzaRatingById(w http.ResponseWriter, r *http.Request) {
	//* establish connection to the database
	db := Connect()
	defer db.Close()

	//* mux.Vars(r) returns a map of the variables in the url, then retrieves the rating from the map
	mux.Vars(r)
	rating := mux.Vars(r)["rating"]

	//* initialize a pizza variable from the pizza struct
	var pizza Pizza

	//* query the database for the pizza with the given ID
	result, err := db.Query("SELECT * FROM Pizza WHERE ID = ?", mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	//* iterate through the result and set the pizza variable to the pizza with the given ID
	for result.Next() {
		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewDecoder(r.Body).Decode(&pizza)

	//* set the pizza rating to the new given rating, after converting the data type from string to float32
	ratingFloat, _ := strconv.ParseFloat(rating, 32)
	pizza.Rating = float32(ratingFloat)

	//* query the database to update the rating of the pizza with the given ID
	_, err = db.Exec("UPDATE Pizza SET Rating = ? WHERE ID = ?", pizza.Rating, pizza.ID)
	if err != nil {
		panic(err.Error())
	}

	log.Println("Pizza rating updated: ", mux.Vars(r)["id"], " to ", pizza.Rating, " stars ", pizza.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizza)
}

// * updatePizzaById updates any field of a pizza with the given ID from the database
func updatePizzaById(w http.ResponseWriter, r *http.Request) {

	//* fetch Name, Price, Picture, Description, Category, Ingredients, Rating from request body
	name := r.FormValue("name")
	price := r.FormValue("price")
	description := r.FormValue("description")
	category := r.FormValue("category")
	ingredients := r.FormValue("ingredients")
	rating := r.FormValue("rating")
	picture := r.FormValue("picture")
	id := r.FormValue("id")

	//* establish connection to the database
	db := Connect()
	defer db.Close()

	//* Check which of the fields are not empty and update the respective fields of the pizza with the given ID
	if name != "" {
		_, err := db.Exec("UPDATE Pizza SET Name = ? WHERE ID = ?", name, id)
		if err != nil {
			panic(err.Error())
		}
	}
	if price != "" {
		floatPrice, _ := strconv.ParseFloat(price, 32)
		_, err := db.Exec("UPDATE Pizza SET Price = ? WHERE ID = ?", floatPrice, id)
		if err != nil {
			panic(err.Error())
		}

	}
	if description != "" {
		_, err := db.Exec("UPDATE Pizza SET Description = ? WHERE ID = ?", description, id)
		if err != nil {
			panic(err.Error())
		}
	}
	if category != "" {
		_, err := db.Exec("UPDATE Pizza SET Category = ? WHERE ID = ?", category, id)
		if err != nil {
			panic(err.Error())
		}
	}
	if ingredients != "" {
		_, err := db.Exec("UPDATE Pizza SET Ingredients = ? WHERE ID = ?", ingredients, id)
		if err != nil {
			panic(err.Error())
		}
	}
	if rating != "" {
		floatRating, _ := strconv.ParseFloat(rating, 32)
		_, err := db.Exec("UPDATE Pizza SET Rating = ? WHERE ID = ?", floatRating, id)
		if err != nil {
			panic(err.Error())
		}
	}
	if picture != "" {
		_, err := db.Exec("UPDATE Pizza SET Picture = ? WHERE ID = ?", picture, id)
		if err != nil {
			panic(err.Error())
		}
	}
}
