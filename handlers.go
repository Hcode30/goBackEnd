package main

import (
	"log"
	"net/http"
	"time"
)



func verifyCookie(r *http.Request) bool {
  token, err := getCookie(r,"token")
  if err != nil {
    log.Println("[ERROR]: ", err)
    return false
  }
  return verifyToken(token)
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
  if !verifyCookie(r) {
    http.Redirect(w, r, "/login/", http.StatusFound)
    return
  }
  renderTemplates(w, "home", nil)
}



func ViewAllHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplates(w, "views", data)
}

func renderTemplates(w http.ResponseWriter, tmpl string, page []Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", page)
	if err != nil {
		log.Printf("[Error] in viewing all pages: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		log.Printf("[Warning]: %s", err)
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", page)
}
func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
  println("body", body)
	author := r.FormValue("author")
	page := &Page{Title: title, Body: []byte(body), Author: author, CreatedAt: time.Now(), ModifiedAt: time.Now()}
	if len(body) != 0 {
		page.save()
	}
	// renderTemplate(w, "edit", page) // don't have multiple responses
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title, CreatedAt: time.Now(), ModifiedAt: time.Now()}
	}
	renderTemplate(w, "edit", page)
}

func MyHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
