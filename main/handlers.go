// ! This file contains the functions(handlers) for the user routes
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hamza-starcevic/pizzaapi/pkg/authPizza"
	db "github.com/hamza-starcevic/pizzaapi/pkg/db"
)

// * Struct used for indicating the status of an operation
type successJson struct {
	Success bool   `json:"success"`
	ID      string `json:"id,omitempty"`
	Message string `json:"message"`
}

type successJsonInt struct {
	Success bool   `json:"success"`
	ID      int    `json:"id,omitempty"`
	Message string `json:"message"`
}

// * fetchPizza returns an array with all pizzas from the database
func fetchPizza(w http.ResponseWriter, r *http.Request) {

	//* establish connection to the database
	db := db.Connect()
	defer db.Close()

	//* initialize an array of pizzas(structs of type Pizza)
	var pizzaArr []Pizza

	//* query the database for all pizzas
	result, err := db.Query("SELECT * FROM Pizza")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error querying the database"})
	}
	defer result.Close()

	//* iterate through the result and append each pizza to the array
	for result.Next() {
		var pizza Pizza
		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error loading the values for pizza"})
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
	db := db.Connect()
	defer db.Close()

	//* mux.Vars(r) returns a map of the variables in the url
	mux.Vars(r)

	//* initialize a pizza(struct of type Pizza)
	var pizza Pizza

	//* query the database for the pizza with the given ID, where the ID is the first variable in the url
	result, err := db.Query("SELECT * FROM Pizza WHERE ID = ?", mux.Vars(r)["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error querying the database"})
	}
	defer result.Close()

	//* iterate through the result and add the values to the pizza variable
	for result.Next() {
		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error loading the values for pizza"})
		}
	}

	//* set the content type to json and encode the pizza
	log.Println("Pizza searched: ", mux.Vars(r)["id"], pizza.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizza)
}

// * createPizza creates a new entry in the database
func createPizza(w http.ResponseWriter, r *http.Request) {

	err := checkAdmin(verifyToken(r))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "You are not authorized to perform this action"})
		return
	}

	//* establish connection to the database
	db := db.Connect()
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
	_, err = db.Exec(`INSERT INTO
							Pizza (ID, Picture, Name, Price, Description, Category, Ingredients, Rating, Date_added)
					   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		pizza.ID, pizza.Picture, pizza.Name, pizza.Price, pizza.Description, pizza.Category, pizza.Ingredients, pizza.Rating, pizza.DateAdded)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error inserting the pizza into the database"})
	}

	//* set the content type to json and encode the pizza struct that was just created
	log.Println("Pizza created: ", pizza.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizza)
}

// * fetchPizzaByCategory returns an array with all the pizzas in the given category
func fetchPizzaByCategory(w http.ResponseWriter, r *http.Request) {

	//* establish connection to the database
	db := db.Connect()
	defer db.Close()

	//* mux.Vars(r) returns a map of the variables in the url
	mux.Vars(r)

	//* initialize an array of pizzas(structs of type Pizza)
	var pizzaArr []Pizza

	//* query the database for all pizzas in the given category, where the category is the first variable in the url
	result, err := db.Query("SELECT * FROM Pizza WHERE Category = ?", mux.Vars(r)["category"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error querying the database"})
	}
	defer result.Close()

	//* iterate through the result and append each pizza to the array
	for result.Next() {
		var pizza Pizza

		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error loading the values for pizza"})
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
	db := db.Connect()
	defer db.Close()

	//* mux.Vars(r) returns a map of the variables in the url, then retrieves the search string from the map
	mux.Vars(r)
	name := mux.Vars(r)["name"]

	//* initialize an array of pizzas(structs of type Pizza)
	var pizzaArr []Pizza

	//* query the database for all pizzas that contain the given string in their name
	result, err := db.Query("SELECT * FROM Pizza WHERE Name LIKE ?", "%"+mux.Vars(r)["name"]+"%")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error querying the database"})
	}
	defer result.Close()

	//* iterate through the result and append each pizza to the array
	for result.Next() {
		var pizza Pizza

		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error loading the values for pizza"})
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

	err := checkAdmin(verifyToken(r))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "You are not authorized to perform this action"})
		return
	}
	//* establish connection to the database
	db := db.Connect()
	defer db.Close()

	//* mux.Vars(r) returns a map of the variables in the url, then retrieves the ID from the map
	mux.Vars(r)

	//* query the database to delete the pizza with the given ID
	result, err := db.Exec("DELETE FROM Pizza WHERE ID = ?", mux.Vars(r)["id"])

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error deleting the pizza from the database"})
	}

	//* set the content type to json and encode successJson struct depending on whether the pizza was deleted or not
	effect, _ := result.RowsAffected()
	w.Header().Set("Content-Type", "application/json")

	if effect == 0 {
		w.WriteHeader(http.StatusNotFound)
		var fail = successJson{Success: false, ID: mux.Vars(r)["id"], Message: "Pizza not found"}
		json.NewEncoder(w).Encode(fail)

	} else {
		w.WriteHeader(http.StatusOK)
		var success = successJson{Success: true, ID: mux.Vars(r)["id"], Message: "Pizza deleted"}
		json.NewEncoder(w).Encode(success)
	}

	log.Println("Pizza Deleted: ", mux.Vars(r)["id"])
}

// * updatePizzaRatingById updates the rating of a pizza with the given ID from the database
func updatePizzaRatingById(w http.ResponseWriter, r *http.Request) {

	err := checkAdmin(verifyToken(r))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "You are not authorized to perform this action"})
		return
	}
	//* establish connection to the database
	db := db.Connect()
	defer db.Close()

	//* mux.Vars(r) returns a map of the variables in the url, then retrieves the rating from the map
	mux.Vars(r)
	rating := mux.Vars(r)["rating"]

	//* initialize a pizza variable from the pizza struct
	var pizza Pizza

	//* query the database for the pizza with the given ID
	result, err := db.Query("SELECT * FROM Pizza WHERE ID = ?", mux.Vars(r)["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error querying the database"})
	}
	defer result.Close()

	//* iterate through the result and set the pizza variable to the pizza with the given ID
	for result.Next() {
		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error updating the rating of the pizza"})
		}
	}

	json.NewDecoder(r.Body).Decode(&pizza)

	//* set the pizza rating to the new given rating, after converting the data type from string to float32
	ratingFloat, _ := strconv.ParseFloat(rating, 32)
	pizza.Rating = float32(ratingFloat)

	//* query the database to update the rating of the pizza with the given ID
	_, err = db.Exec("UPDATE Pizza SET Rating = ? WHERE ID = ?", pizza.Rating, pizza.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error updating the rating of the pizza"})
	}

	log.Println("Pizza rating updated: ", mux.Vars(r)["id"], " to ", pizza.Rating, " stars ", pizza.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizza)
}

// * updatePizzaById updates any field of a pizza with the given ID from the database
func updatePizzaById(w http.ResponseWriter, r *http.Request) {

	err := checkAdmin(verifyToken(r))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "You are not authorized to perform this action"})
		return
	}
	var pizza Pizza
	json.NewDecoder(r.Body).Decode(&pizza)

	//* fetch Name, Price, Picture, Description, Category, Ingredients, Rating from request body
	name := pizza.Name
	price := pizza.Price
	description := pizza.Description
	category := pizza.Category
	ingredients := pizza.Ingredients
	rating := pizza.Rating
	picture := pizza.Picture
	id := mux.Vars(r)["id"]

	//* establish connection to the database
	db := db.Connect()
	defer db.Close()

	//* query the database to update the pizza with the given ID
	_, err = db.Exec("UPDATE Pizza SET Name = ?, Price = ?, Description = ?, Category = ?, Ingredients = ?, Rating = ?, Picture = ? WHERE ID = ?", name, price, description, category, ingredients, rating, picture, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error updating the pizza"})
	}
	pizza.ID = id
	//fetch the DateAdded field from the database
	var dateAdded string
	_ = db.QueryRow("SELECT Date_added FROM Pizza WHERE ID = ?", id).Scan(&dateAdded)
	//Log the data type of dateAdded
	pizza.DateAdded = dateAdded
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizza)
}

// * createOrder creates a new entry in the Orders table in the database,
// * as well as a new entry in the OrderItems table for each item in the order
func createOrder(w http.ResponseWriter, r *http.Request) {

	//* fetch the user ID from the token
	token, err := verifyToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "You are not authorized to perform this action"})
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	userEmail := claims["email"].(string)

	//* fetch TotalPrice from the request body
	var order Order
	json.NewDecoder(r.Body).Decode(&order)

	//query the database to find the user with the given ID
	db := db.Connect()
	defer db.Close()
	log.Println("This is the user ID: ", userID)
	log.Println("This is the user email: ", userEmail)
	log.Println("This is the orice : ", order.Total_price)
	log.Println("This is the City : ", order.City)
	log.Println("This is the Street : ", order.Street)
	log.Println("This is the zip : ", order.Zipcode)

	//* query the database to insert the order into the database where the user id is a foreign key from the Users table and retrieve the ID of the order

	result, err := db.Exec("INSERT INTO PizzaOrder (Total_price, UserID, UserEmail, City, Zipcode, Street) VALUES (?, ?, ?, ?, ?, ?)",
		order.Total_price,
		userID,
		userEmail,
		order.City,
		order.Zipcode,
		order.Street)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error creating the order"})
		return
	}
	orderID, _ := result.LastInsertId()
	log.Println("This is the order ID: ", orderID)
	log.Println("I can get to this point")

	//* fetch the ID of the order that was just created

	//* range through order.PizzaIDS and query the database to insert the pizzaID and orderId into the order_items table
	for _, pizzaID := range order.PizzaIDS {
		//Find the name of the pizza with the given ID
		var pizzaName string
		_ = db.QueryRow("SELECT Name FROM Pizza WHERE ID = ?", pizzaID).Scan(&pizzaName)

		_, err = db.Exec("INSERT INTO order_items (item_id, order_ID, item_name) VALUES (?, ?, ?)", pizzaID, orderID, pizzaName)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error creating the order_items"})
			return
		}
	}
	order.UserEmail = userEmail
	order.ID = int64(orderID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// *getOrders returns all the orders in the database
func getOrders(w http.ResponseWriter, r *http.Request) {

	err := checkAdmin(verifyToken(r))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "You are not authorized to perform this action"})
		return
	}

	var orders []Order
	//* establish connection to the database
	db := db.Connect()
	defer db.Close()

	//* query the database to get all the orders
	rows, err := db.Query("SELECT * FROM PizzaOrder")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error fetching the orders from the database"})
		return
	}
	defer rows.Close()

	log.Println("I can get to this point")
	//* range through the rows and append the orders to the orders slice

	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.Total_price, &order.UserID, &order.UserEmail, &order.City, &order.Zipcode, &order.Street, &order.Date_added)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error loading the orders"})
			return
		}
		orders = append(orders, order)
	}

	// * range through the orders slice and query the database to get the items in each order
	for i, order := range orders {
		//query the database to get the items in the order
		rows, err := db.Query("SELECT item_id FROM order_items WHERE order_ID = ?", order.ID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error fetching the items from the database"})
			return
		}
		defer rows.Close()

		//* range through the rows and append the pizzaIDs to the orders slice

		for rows.Next() {
			var pizzaID string
			err := rows.Scan(&pizzaID)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error loading the items"})
				return
			}

			orders[i].PizzaIDS = append(orders[i].PizzaIDS, pizzaID)
		}
	}

	//* range through the pizzaIDS and query the database to append all of the Pizza.Picture to the order.Pizzadetails slice
	for i, order := range orders {
		for _, pizzaID := range order.PizzaIDS {
			var pizza Pizza
			_ = db.QueryRow("SELECT Picture FROM Pizza WHERE ID = ?", pizzaID).Scan(&pizza.Picture)
			orders[i].PizzaDetails = append(orders[i].PizzaDetails, pizza.Picture)
		}

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func getOrdersbyId(w http.ResponseWriter, r *http.Request) {

	token, err := verifyToken(r)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "You are not authorized to perform this action"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	var orders []Order
	//* establish connection to the database
	db := db.Connect()
	defer db.Close()

	log.Println("This is the user ID: ", userID)
	//* query the database to get all the orders
	rows, err := db.Query("SELECT * FROM PizzaOrder WHERE UserID = ?", userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error fetching the orders from the database"})
		return
	}
	defer rows.Close()

	log.Println("I can get to this point")
	//* range through the rows and append the orders to the orders slice

	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.Total_price, &order.UserID, &order.UserEmail, &order.City, &order.Zipcode, &order.Street, &order.Date_added)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error loading the orders"})
			return
		}
		orders = append(orders, order)
	}

	// * range through the orders slice and query the database to get the items in each order
	for i, order := range orders {
		//query the database to get the items in the order
		rows, err := db.Query("SELECT item_id FROM order_items WHERE order_ID = ?", order.ID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error fetching the items from the database"})
			return
		}
		defer rows.Close()

		//* range through the rows and append the pizzaIDs to the orders slice

		for rows.Next() {
			var pizzaID string
			err := rows.Scan(&pizzaID)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error loading the items"})
				return
			}

			orders[i].PizzaIDS = append(orders[i].PizzaIDS, pizzaID)
		}
	}

	//* range through the pizzaIDS and query the database to append all of the pizzas to the order.Pizzadetails slice
	for i, order := range orders {
		for _, pizzaID := range order.PizzaIDS {
			var pizza Pizza
			err := db.QueryRow("SELECT * FROM Pizza WHERE ID = ?", pizzaID).Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error fetching the pizzas from the database"})
				return
			}
			orders[i].PizzaDetails = append(orders[i].PizzaDetails, pizza.Picture)
		}

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// * deleteOrder deletes an order from the database
func deleteOrderById(w http.ResponseWriter, r *http.Request) {

	err := checkAdmin(verifyToken(r))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "You are not authorized to perform this action"})
		return
	}

	//* fetch the ID of the order from the request body
	id := mux.Vars(r)["id"]

	//* query the database to delete the order with the given ID
	db := db.Connect()
	defer db.Close()

	// * Check if the order exists
	var order Order
	err = db.QueryRow("SELECT ID FROM PizzaOrder WHERE id = ?", id).Scan(&order.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "The order does not exist"})
		return
	}
	_, err = db.Exec("DELETE FROM PizzaOrder WHERE ID = ?", id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error deleting the order"})
		return
	}

	//* query the database to delete the order_items with the given order ID
	_, err = db.Exec("DELETE FROM order_items WHERE order_ID = ?", id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error deleting the order_items"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error converting the id to an int"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(successJsonInt{Success: true, ID: idInt, Message: "Order deleted successfully"})
}

// * fetchUsers fetches all the users from the database
func fetchUsers(w http.ResponseWriter, r *http.Request) {

	err := checkAdmin(verifyToken(r))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "You are not authorized to perform this action"})
		return
	}
	//* establish connection to the database
	db := db.Connect()
	defer db.Close()

	//* query the database to fetch all the users
	rows, err := db.Query("SELECT * FROM Users")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error fetching the users"})
	}

	var users []authPizza.User
	var blank string
	//* range through the rows and append them to the fetchedusers array
	for rows.Next() {
		var user authPizza.User
		err = rows.Scan(&user.Email, &blank, &user.ID, &user.Role, &user.Date_added)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(successJson{Success: false, Message: "Error fetching the users"})
		}
		//Dont show the password
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// * deleteUser deletes a user corresponding to the id in the url from the database
func deleteUserById(w http.ResponseWriter, r *http.Request) {

	err := checkAdmin(verifyToken(r))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "You are not authorized to perform this action"})
		return
	}
	//* fetch the user ID from the url
	userID := mux.Vars(r)["id"]

	//* establish connection to the database
	db := db.Connect()
	defer db.Close()

	//* query the database to delete the user with the given ID
	_, err = db.Exec("DELETE FROM Users WHERE ID = ?", userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, ID: userID, Message: "Error deleting the user"})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(successJson{Success: true, ID: userID, Message: "User deleted successfully"})
}

// * updateUser updates a users role and/or email by the id, in the database
func updateUserById(w http.ResponseWriter, r *http.Request) {

	err := checkAdmin(verifyToken(r))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "You are not authorized to perform this action"})
		return
	}
	//* fetch the user ID from the url
	userID := mux.Vars(r)["id"]

	//* fetch the user from the request body
	var user authPizza.User
	json.NewDecoder(r.Body).Decode(&user)

	//* establish connection to the database
	db := db.Connect()
	defer db.Close()

	//* query the database to update the user with the given ID
	_, err = db.Exec("UPDATE Users SET Email = ?, Role = ? WHERE ID = ?", user.Email, user.Role, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, ID: userID, Message: "Error updating the user"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	log.Println(user.Email)
	log.Println(user.Role)
	log.Println(userID)
	json.NewEncoder(w).Encode(successJson{Success: true, ID: userID, Message: "User updated successfully"})
}

// * verifyToken verifies the token in the request header, and returns it if it is valid
func verifyToken(r *http.Request) (*jwt.Token, error) {
	//Check for bearer token
	tokenString := r.Header.Get("Authorization")
	//Parse the token

	if tokenString == "" {
		return nil, errors.New("token not found")
	}
	//strip the bearer part of the token
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	if tokenString != "" {
		//Parse HS256 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

				return nil, fmt.Errorf("there was an error")
			}

			return []byte("SecretYouShouldHide"), nil
		})
		if err != nil {
			log.Println("Error parsing token: ", err)
			return nil, err
		}
		// if err != nil {
		// 	log.Println("Error parsing token not ok this is the token: ", token)

		// 	return nil, err
		// }
		return token, nil
	}
	return nil, nil
}

// * checkAdmin checks if the user is an admin
func checkAdmin(t *jwt.Token, err error) error {
	if err != nil {
		log.Println("Error parsing token not ok this is the token: ", t)
		return err
	}

	claims := t.Claims.(jwt.MapClaims)
	if claims["isAdmin"] == false {
		return errors.New("not an admin")
	}

	return nil
}

// * verifyToken verifies the token in the request header, and returns it if it is valid
func verifyTokenTest(w http.ResponseWriter, r *http.Request) {
	//Check for bearer token
	tokenString := mux.Vars(r)["token"]
	//Parse the token

	if tokenString == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(successJson{Success: false, Message: "No token found"})
		return
	}
	if tokenString != "" {
		//Parse HS256 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

				return nil, fmt.Errorf("there was an error")
			}

			return []byte("SecretYouShouldHide"), nil
		})
		if err != nil {
			log.Println("Error parsing token: ", err)
		}
		// if err != nil {
		// 	log.Println("Error parsing token not ok this is the token: ", token)

		// 	return nil, err
		// }

		claims := token.Claims.(jwt.MapClaims)
		log.Println(claims["email"])
	}
}
