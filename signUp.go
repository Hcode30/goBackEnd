package main

import (
	"log"
	"net/http"
)

var users = []User{}

// redirect to viewAll page if signed up successfully
func signUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		log.Println("email: ", email, "password: ", password, "firstName: ", firstName, "lastName: ", lastName)
		user, err := MakeUser(email, firstName, lastName, password)
		if err != nil {
      log.Println("[ERROR]: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, *user)
		log.Println("[INFO]: Signed Up successfully")
		log.Println("User: ", user.UserName, "Email: ", user.Email,
			"Password: ", password)
		http.Redirect(w, r, "/viewAll/", http.StatusFound)
		// }
		renderTemplate(w, "signup", nil)
	}
}
func signUpPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "signup.html", nil)
	}
}
