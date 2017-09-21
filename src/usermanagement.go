package main

import (
	"database/sql"
	"fmt"
)

// Translate POST requests into more basic parameters
// and pass to specific function
func RequestTranslator(w http.ResponseWriter, r *http.Request) {

}

// Check username and password, pass back secure cookie
func Login() {

}

// Called to verify cookie token
func VerifyLogin() {

}

// Unregister cookie - clear cached token from database
func Logout() {

}

// Insert new user into database
func CreateUser() {

}
