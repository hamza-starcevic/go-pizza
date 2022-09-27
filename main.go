package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var pizzaArr = []Pizza{
	{ID: 1, Picture: "https://st.depositphotos.com/1003814/5052/i/950/depositphotos_50523105-stock-photo-pizza-with-tomatoes.jpg", Name: "Margherita", Price: 10, Description: "THand tossed bbq rib meat lovers pan Chicago style bianca meatball, party garlic sauce. Buffalo sauce beef marinara bbq rib hand tossed pineapple. Mushrooms white garlic platter, pepperoni fresh tomatoes ricotta mayo. Pizza bacon & tomato bbq mushrooms extra sauce. Pizza roll bianca banana peppers, extra cheese Chicago style personal buffalo sauce extra sauce pan pepperoni melted cheese sausage sausage lasagna meatball. Lasagna mushrooms pineapple mozzarella, red onions ricotta ranch burnt mouth.", Category: "Veg"},
	{ID: 2, Picture: "https://pbs.twimg.com/media/FTNbp3WXEAAIZcb.jpg", Name: "Pepperoni", Price: 15, Description: "Sauteed onions lasagna parmesan steak bacon & tomato thin crust pan chicken wing extra sauce pie spinach beef broccoli NY style marinara. Thin crust white garlic anchovies, platter chicken wing peppers party white pizza onions lasagna. Red onions ranch pizza, garlic parmesan green bell peppers fresh tomatoes thin crust philly chicken. Thin crust large spinach pepperoni bacon, chicken mozzarella garlic platter thin crust personal parmesan meatball. Bbq rib meatball beef anchovies mayo pepperoni bianca.", Category: "Non-Veg"},
	{ID: 3, Picture: "https://thumbs.dreamstime.com/b/pizza-rustic-italian-mozzarella-cheese-basil-leaves-35669930.jpg", Name: "Veggie", Price: 20, Description: "THawaiian sauteed onions black olives onions bbq rib spinach garlic mayo. Onions fresh tomatoes marinara, bbq sauce pepperoni pie black olives pork lasagna chicken wing. Spinach meatball peppers, philly chicken meatball pork anchovies mayo. Steak platter pizza roll white garlic, string cheese chicken wing mayo", Category: "Veg"},
	{ID: 4, Picture: "https://pizzafellaz.ca/wp-content/uploads/2021/12/food-pizza-still-life-hd-wallpaper-preview.jpg", Name: "Chicken Sausage", Price: 25, Description: "NY style pepperoni buffalo sauce, thin crust philly chicken chicken ham hawaiian pizza personal hand tossed pie banana peppers. Bbq platter bbq sauce, NY style philly chicken personal garlic sauce parmesan onions meatball peppers meat lovers green bell peppers stuffed crust. Garlic platter banana peppers green bell peppers bacon extra sauce pizza. Ranch personal buffalo sauce ricotta pan mozzarella meatball, pork platter pepperoni sauteed onions burnt mouth Chicago style spinach. String cheese thin crust ham, bianca marinara chicken wing meatball pizza red onions white pizza black olives bbq rib banana peppers anchovies.", Category: "Non-Veg"},
	{ID: 5, Picture: "https://images5.alphacoders.com/399/thumb-1920-399563.jpg", Name: "Chicken Tikka", Price: 30, Description: "Anchovies chicken hawaiian black olives. Large bianca chicken fresh tomatoes mushrooms, green bell peppers onions melted cheese. Black olives extra sauce bbq rib, bacon & tomato garlic sausage beef philly chicken mushrooms stuffed sausage bianca ricotta. Steak party bbq rib thin crust personal chicken wing sausage melted cheese ranch onions mozzarella. Deep crust pan red onions chicken melted cheese extra sauce.", Category: "Non-Veg"},
}

const port = ":8070"

func main() {
	//Import pizzaArr to the database
	// Init()
	// for _, pizza := range pizzaArr {
	// 	_, err := db.Exec("INSERT INTO pizza (picture, name, price, description, category) VALUES ($1, $2, $3, $4, $5)", pizza.Picture, pizza.Name, pizza.Price, pizza.Description, pizza.Category)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	//Creating Gorilla router
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	r.HandleFunc("/pizza", fetchPizza).Methods("GET")
	r.HandleFunc("/pizzaByCategory/{category}", fetchPizzaByCategory).Methods("GET")
	r.HandleFunc("/searchPizza/{searchTerm}", searchPizza).Methods("GET")
	//Starting server
	log.Printf("Server started on port %s", port)
	http.ListenAndServe(port, r)

}
