package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type successJson struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

func fetchPizza(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	var pizzaArr []Pizza
	result, err := db.Query("SELECT * FROM Pizza")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var pizza Pizza
		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			panic(err.Error())
		}
		pizzaArr = append(pizzaArr, pizza)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizzaArr)
}

func fetchPizzaById(w http.ResponseWriter, r *http.Request) {
	//get pizza by id from db
	db := Connect()
	mux.Vars(r)

	var pizza Pizza
	result, err := db.Query("SELECT * FROM Pizza WHERE ID = ?", mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			panic(err.Error())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizza)
}

func createPizza(w http.ResponseWriter, r *http.Request) {
	//create pizza in db
	db := Connect()
	var pizza Pizza
	json.NewDecoder(r.Body).Decode(&pizza)
	pizza.DateAdded = time.Now().Format("2006-01-02 15:04:05")
	print(pizza.Rating)
	result, err := db.Exec("INSERT INTO Pizza (Picture, Name, Price, Description, Category, Ingredients, Rating, Date_added) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", pizza.Picture, pizza.Name, pizza.Price, pizza.Description, pizza.Category, pizza.Ingredients, pizza.Rating, pizza.DateAdded)
	if err != nil {
		panic(err.Error())
	}
	ID, _ := result.LastInsertId()

	pizza.ID = int(ID)

	print(pizza.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizza)
}

func fetchPizzaByCategory(w http.ResponseWriter, r *http.Request) {
	//get pizza by category from db
	db := Connect()
	mux.Vars(r)

	var pizzaArr []Pizza
	result, err := db.Query("SELECT * FROM Pizza WHERE Category = ?", mux.Vars(r)["category"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var pizza Pizza
		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			panic(err.Error())
		}
		pizzaArr = append(pizzaArr, pizza)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizzaArr)
}

func searchPizza(w http.ResponseWriter, r *http.Request) {
	//search pizza by name from db
	db := Connect()
	mux.Vars(r)
	name := mux.Vars(r)["name"]
	print(name)
	var pizzaArr []Pizza
	result, err := db.Query("SELECT * FROM Pizza WHERE Name LIKE ?", "%"+mux.Vars(r)["name"]+"%")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var pizza Pizza
		err := result.Scan(&pizza.ID, &pizza.Picture, &pizza.Name, &pizza.Price, &pizza.Description, &pizza.Category, &pizza.Ingredients, &pizza.Rating, &pizza.DateAdded)
		if err != nil {
			panic(err.Error())
		}
		pizzaArr = append(pizzaArr, pizza)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pizzaArr)
}

func deletePizzaById(w http.ResponseWriter, r *http.Request) {
	//delete pizza by id from db
	db := Connect()
	mux.Vars(r)

	result, err := db.Exec("DELETE FROM Pizza WHERE ID = ?", mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("result: %v\n", result)
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

	print(result)
}
