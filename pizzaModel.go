package main

type Pizza struct {
	ID          int     `json:"id"`
	Picture     string  `json:"picture"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Ingredients string  `json:"ingredients"`
	Rating      float32 `json:"rating"`
	DateAdded   string  `json:"date_added"`
}
