package main

// * The struct used to represent a pizza
// * The struct tags are used to specify the column names in the database
// * The data types are used to specify the data types in the database
// * The json tags in backticks represent what the data will look like once parsed to json
type Pizza struct {
	ID          string  `json:"id"`
	Picture     string  `json:"picture"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Ingredients string  `json:"ingredients"`
	Rating      float32 `json:"rating"`
	DateAdded   string  `json:"date_added"`
}
