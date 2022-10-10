package authPizza

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	db "github.com/hamza-starcevic/pizzaapi/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

type successJson struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
	Message string `json:"message"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// * Data is parsed from json
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//* Check if the user already exists
	db := db.Connect()
	defer db.Close()

	var dbUser User
	err = db.QueryRow(`SELECT
						ID, email, password
					FROM
						Users
					WHERE
						email = ?`,
		user.Email).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, ID: "", Message: "User already exists"})
		return
	}
	// Create uuid
	ID := uuid.New().String()
	user.ID = ID
	user.Password = string(hashedPassword)

	// * The user is inserted into the database
	_, err = db.Exec(`INSERT INTO
						Users (email, password, ID)
					VALUES (?, ?, ?)`,
		user.Email, string(hashedPassword), user.ID)
	if err != nil {
		log.Println("Error inserting user into database: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// * Take out the password from the response

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
