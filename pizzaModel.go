package main

type Pizza struct {
	ID          int    `json:"id"`
	Picture     string `json:"picture"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Category    string `json:"category"`
}
