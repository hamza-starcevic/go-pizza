// ! This file contains the function to initialize
// ! the database connection using a specific sql driver
package db

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// * Connect initializes the database connection
func Connect() *sql.DB {
	//! CRUCIAL: The database must be created beforehand
	//! The table must be created beforehand, and named Pizza
	//! The table must have the following columns:
	//! ID varchar (primary key) NOT NULL,
	//! Picture varchar NOT NULL,
	//! Name varchar(255) NOT NULL,
	//! Category varchar(255) NOT NULL,
	//! Price float NOT NULL,
	//! Rating float NOT NULL,
	//! Description varchar(255) NOT NULL
	//! Date_added timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP

	//* The database connection is initialized
	//! CRUCIAL: The connection string is unique to your database
	//! exapmple dns connection string: "username:password@tcp(address:port)/databaseName"
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	dns := os.Getenv("DNS")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err.Error())
	}

	return db
}
