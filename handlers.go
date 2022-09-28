package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
)

func fetchPizza(w http.ResponseWriter, r *http.Request) {
	//fetch all pizzas to pizzaArray from firestore database
	client, ctx := Init()
	pizzaArray := []Pizza{}
	//Get a document with the query
	iter := client.Collection("pizzas").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		var pizza Pizza
		doc.DataTo(&pizza)
		pizzaArray = append(pizzaArray, pizza)
	}
	// Converting array to formatted JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pizzaArray)
}

// get handler that takes a paramater of category and returns the pizza of that category
func fetchPizzaByCategory(w http.ResponseWriter, r *http.Request) {
	// Fetching category from URL
	vars := mux.Vars(r)
	category := vars["category"]

	client, ctx := Init()
	pizzaArray := []Pizza{}
	//Get a document with the query
	iter := client.Collection("pizzas").Where("Category", "==", category).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		var pizza Pizza
		doc.DataTo(&pizza)
		pizzaArray = append(pizzaArray, pizza)
	}
	// Converting array to formatted JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pizzaArray)
}

func fetchPizzaById(w http.ResponseWriter, r *http.Request) {
	//fetch pizza from firestore database by id
	client, ctx := Init()

	vars := mux.Vars(r)
	id := vars["id"]

	//Get a document with the query
	doc := client.Collection("pizzas").Where("PID", "==", id).Documents(ctx)
	//Get the document
	docSnap, err := doc.Next()
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
	}
	//Get the data from the document
	var pizza Pizza
	docSnap.DataTo(&pizza)
	// Converting array to formatted JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pizza)

}

// get handler that takes a search term and returns the pizza that matches the search term
func searchPizza(w http.ResponseWriter, r *http.Request) {
	// Fetching search term from URL
	vars := mux.Vars(r)
	searchTerm := vars["searchTerm"]
	//check is search term is empty
	if searchTerm == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pizzaArr)
		return
	}
	print(searchTerm)
	sortedPizza := []Pizza{}
	// Filtering pizza by search term
	client, ctx := Init()
	iter := client.Collection("pizzas").Where("Name", "==", searchTerm).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		var pizza Pizza
		doc.DataTo(&pizza)
		sortedPizza = append(sortedPizza, pizza)
	}
	// Converting array to formatted JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sortedPizza)
}

// post handler that takes a pizza and adds it to the database
func addPizza(w http.ResponseWriter, r *http.Request) {
	print("hit")

	//Count the number of pizzas in the database
	//Get a document with the query
	client, ctx := Init()

	doc, err := client.Collection("pizzas").Documents(ctx).GetAll()
	//get number of documents
	numOfPizzas := len(doc)
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
	}
	print(numOfPizzas)

	//fetch json from post and declare variables with their values

	var pizza Pizza
	_ = json.NewDecoder(r.Body).Decode(&pizza)
	//add pizza to firestore database
	pizza.PID = strconv.Itoa(numOfPizzas + 1)
	_, _, err = client.Collection("pizzas").Add(ctx, pizza)
	print(pizza.PID)

	if err != nil {
		log.Fatalf("Failed to add pizza to firestore: %v", err)
	}
	//add pizza to pizzaArr
	pizzaArr = append(pizzaArr, pizza)
	// Converting array to formatted JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pizzaArr)

}
