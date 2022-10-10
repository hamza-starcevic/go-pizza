package authPizza

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	db "github.com/hamza-starcevic/pizzaapi/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

type jwtJson struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Token   string `json:"token"`
	IsAdmin bool   `json:"isAdmin"`
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	// * Data is parsed from json
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(user.Email)
	log.Println(user.Password)
	log.Println(user.Role)
	log.Println(user.ID)
	// * The user is fetched from the database
	db := db.Connect()
	defer db.Close()

	var dbUser User
	err = db.QueryRow(`SELECT
						ID, Email, Password, Role
					FROM
						Users
					WHERE
						email = ?`,
		user.Email).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password, &dbUser.Role)
	log.Println("DATABASE DATA DOWN, LOGIN DATA UP")
	log.Println(dbUser.Email)
	log.Println(dbUser.Password)
	log.Println(dbUser.Role)
	log.Println(dbUser.ID)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, ID: "", Message: "Invalid user credentials"})
		return
	}

	// * The password is compared
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		log.Println("This is the error: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, ID: "", Message: "Invalid user credentials"})
		return
	}

	// * The token is created
	sampleSecretKey := []byte("SecretYouShouldHide")

	token := jwt.New(jwt.SigningMethodHS256)

	log.Println(token, " This is the token")

	// * Check if the user is admin
	var isAdmin bool
	var tempRole string
	err = db.QueryRow(`SELECT
							Role
						FROM
							Users
						WHERE
							email = ?`,
		user.Email).Scan(&tempRole)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(successJson{Success: false, ID: "", Message: "Error checking user role"})
		return
	}
	if tempRole == "admin" {
		isAdmin = true
	} else {
		isAdmin = false
	}
	//Set standard claims expiration token to 1 hour

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = dbUser.ID
	claims["email"] = dbUser.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	if isAdmin {
		claims["isAdmin"] = true
	} else {
		claims["isAdmin"] = false
	}

	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		log.Println("Error while signing the token")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// * The token is sent back to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if isAdmin {
		json.NewEncoder(w).Encode(jwtJson{Token: tokenString, IsAdmin: true, Email: user.Email, ID: dbUser.ID})
	} else {
		json.NewEncoder(w).Encode(jwtJson{Token: tokenString, IsAdmin: false, Email: user.Email, ID: dbUser.ID})
	}
}
