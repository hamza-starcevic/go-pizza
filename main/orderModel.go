package main

type Order struct {
	ID           int64    `json:"id,omitempty"`
	UserID       string   `json:"user_id,omitempty"`
	UserEmail    string   `json:"user_email,omitempty"`
	Total_price  float32  `json:"total_price,omitempty"`
	PizzaIDS     []string `json:"pizza_ids,omitempty"`
	City         string   `json:"city,omitempty"`
	Street       string   `json:"street,omitempty"`
	Zipcode      string   `json:"zipcode,omitempty"`
	Date_added   string   `json:"date_added,omitempty"`
	PizzaDetails []string `json:"pizza_details,omitempty"`
}

type OrderItem struct {
	OrderID   string `json:"order_id"`
	PizzaID   string `json:"pizza_id"`
	PizzaName string `json:"pizza_name"`
}
