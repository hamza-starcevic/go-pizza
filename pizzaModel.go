package main

type Pizza struct {
	PID         string  `json:"id"`
	Picture     string  `json:"picture"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}
