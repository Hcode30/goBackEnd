package main

import (
	"errors"
	"log"
	"net/http"
)

// redirect to login page if not logged in successfully
// redirect to viewAll page if logged in successfully
func logInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")
		user, err := GetUser(email)
		if err != nil {
      log.Println("[ERROR]: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = verifyPassword(user, password)
		if err != nil {
      log.Println("[ERROR]: ", err)
      // reload the login page with the email typed
      http.Redirect(w, r, "/login/", http.StatusNotFound)
			return
		}
    log.Println("[INFO]: Logged in successfully")
		http.Redirect(w, r, "/viewAll/", http.StatusFound)
	}
	renderTemplate(w, "login", nil)
}

func verifyPassword(user *User, password string) error {
	hsh := Argon2idHash{
		time:    3,
		saltLen: 16,
		memory:  12288,
		threads: 1,
		keyLen:  32,
	}
	hashSalt, err := hsh.GenerateHash([]byte(password), user.Password.salt)
	if err != nil {
		return err
	}
	if string(hashSalt.Hash) != string(user.Password.hash) {
		return errors.New("Invalid password")
	}
	return nil
}

func logInPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "login.html", nil)
	}
}
