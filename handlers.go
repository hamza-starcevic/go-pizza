package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func fetchPizza(w http.ResponseWriter, r *http.Request) {
	// Converting array to formatted JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pizzaArr)
}

// get handler that takes a paramater of category and returns the pizza of that category
func fetchPizzaByCategory(w http.ResponseWriter, r *http.Request) {

	// Fetching category from URL
	vars := mux.Vars(r)
	category := vars["category"]
	print(category)
	sortedPizza := []Pizza{}
	// Filtering pizza by category
	for _, pizza := range pizzaArr {
		if pizza.Category == category {
			sortedPizza = append(sortedPizza, pizza)
		}
	}
	// Converting array to formatted JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sortedPizza)
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
	for _, pizza := range pizzaArr {
		if strings.Contains(strings.ToLower(pizza.Name), strings.ToLower(searchTerm)) {
			sortedPizza = append(sortedPizza, pizza)
		}
	}
	// Converting array to formatted JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sortedPizza)
}
