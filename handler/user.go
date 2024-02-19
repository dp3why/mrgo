package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/dp3why/mrgo/constants"
	"github.com/dp3why/mrgo/model"
	"github.com/dp3why/mrgo/service"
	jwt "github.com/form3tech-oss/jwt-go"
)

var mySigningKey = []byte(constants.JWT_SECRET)

func signHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one sign request")
	w.Header().Set("Content-Type", "text/plain")

	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode user data from client %v.\n", err)
		return
	}
	success, err := service.CheckUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, "Failed to read user to ElasticSearch", http.StatusInternalServerError)
		fmt.Printf("Failed to read user from ElasticSearch %v.\n", err)
		return
	}
	if !success {
		http.Error(w, "User does not exist / wrong password", http.StatusBadRequest)
		fmt.Printf("User does not exist / wrong password.\n")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		fmt.Printf("Failed to generate token %v.\n", err)
		return
	}
	w.Write([]byte(tokenString))
}



func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signup request")
	w.Header().Set("Content-Type", "text/plain")

	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode user data from client %v.\n", err)
		return
	}
	if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		fmt.Printf("Invalid username or password.\n")
		return
	}
	success, err := service.AddUser(&user)
	if err != nil {
		http.Error(w, "Failed to save user to ElasticSearch", http.StatusInternalServerError)
		fmt.Printf("Failed to save user to ElasticSearch %v.\n", err)
		return
	}
	if !success {
		http.Error(w, "User already exists", http.StatusBadRequest)
		fmt.Printf("User already exists.\n")
		return
	}
	 
}
