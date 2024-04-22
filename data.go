package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
)

const (
	dataDir = "data/"
)

type Page struct {
	Title string
	Body  []byte
}

var (
	templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
	data      = loadAllData(dataDir)
	dataMutex sync.Mutex
)

func (p *Page) save() {
	dataMutex.Lock()
	defer dataMutex.Unlock()
	for i, page := range data {
		if page.Title == p.Title {
			data[i].Body = p.Body
			return
		}
	}
	data = append(data, *p)
}

func saveAllData(data []Page) bool {
	dataMutex.Lock()
	defer dataMutex.Unlock()
	println("[INFO]: Saving", len(data), "pages...")
	for _, page := range data {
		filepath := dataDir + page.Title + ".txt"
		println("[INFO]: saving", filepath)
		err := os.WriteFile(filepath, page.Body, 0600)
		if err != nil {
			return false
		}
	}
	return true
}

func loadAllData(dataDir string) []Page {
	dataMutex.Lock()
	defer dataMutex.Unlock()
	f, err := os.Open(dataDir)
	if err != nil {
		log.Printf("[Error] opening %s : %s", dataDir, err)
		panic(err)
	}
	files, err := f.Readdir(0)
	if err != nil {
		log.Printf("[Error] opening files in %s : %s", dataDir, err)
		panic(err)
	}
	data := make([]Page, 0, len(files))
	for _, v := range files {
		if !v.IsDir() {
			page := Page{Title: v.Name()[:len(v.Name())-4]}
			body, err := os.ReadFile(dataDir + "/" + v.Name())
			if err != nil {
				log.Printf("[Error] in %s: %s", page.Title, err)
			}
			page.Body = body
			data = append(data, page)
		}
	}
	println("[INFO]: Finished loading data with", len(data), "pages")
	return data
}

func loadPage(title string) (*Page, error) {
	for _, page := range data {
		if page.Title == title {
			return &page, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("can't find page %s!", title))
}
func renderTemplate(w http.ResponseWriter, tmpl string, page *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", page)
	if err != nil {
		log.Printf("[Error] in %s: %s", page.Title, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
