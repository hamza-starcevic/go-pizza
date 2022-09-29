package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "hstarcevic:kenansin@tcp(localhost)/pizzaDB")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Success!")
	return db
}
