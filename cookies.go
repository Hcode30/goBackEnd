package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func addCookie(w http.ResponseWriter, name, value, path string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Secure:   true,
		MaxAge:   0,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func getCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return "", err
	}
	return cookie.Value, nil
}

func generateToken(user User, t time.Time) string {
	// valid for 7 day
	var duratuin int64 = 604800
	token := uint64(user.ID) + uint64(t.Unix()+duratuin)
	strToken := strconv.FormatUint(token, 10)
	formatedID := strconv.FormatUint(user.ID, 10)
	strToken = formatedID + strToken
	// println("new token: ", strToken,"-- len: ",len(strToken))
	return strToken
}

func verifyToken(token string) bool {
	id := token[:18]
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return false
	}
	token = token[18:]
	tokenUint, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return false
	}
	to := idInt + uint64(time.Now().Unix())
	if tokenUint < to {
		log.Println("[WARNING]: Token expired")
		return false
	}
	return true
}
