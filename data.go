package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	dataDir = "data/"
)

type Page struct {
	Title      string
	Body       []byte
	Author     string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

var (
	templates = template.Must(template.ParseFiles(
  // "templates/edit.html",
  // "templates/view.html",
  // "templates/views.html",
  "templates/home.html",
  "templates/signup.html",
  "templates/components.html",
  "templates/login.html"))
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
			data[i].Author = p.Author
			return
		}
	}
	data = append(data, *p)
}

func saveAllData(data []Page) bool {
	if len(data) == 0 {
		println("[INFO]: No data to save")
		return false
	}
	dataMutex.Lock()
	defer dataMutex.Unlock()
	println("[INFO]: Saving", len(data), "pages...")
	for _, page := range data {
		filepath := dataDir + page.Title + ".txt"
		println("[INFO]: saving", filepath)
		data := fmt.Sprintf("[%s]\n[%s]\n%s", page.Author, page.CreatedAt.Format("2006-01-02"), page.Body)
		err := os.WriteFile(filepath, []byte(data), 0600)
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
			author, dateCreated, displayedData, err := extractMetadata(string(body))
			if err != nil {
				log.Printf("[Error]: %s", err)
			}
			page.CreatedAt = dateCreated
			page.ModifiedAt = v.ModTime()
			page.Author = author
			page.Body = []byte(displayedData)
			data = append(data, page)
		}
	}
	if len(data) == 0 {
		println("[INFO]: No data found")
	} else {
		println("[INFO]: Finished loading data with", len(data), "pages")
	}
	return data
}
func extractMetadata(body string) (string, time.Time, string, error) {
    lines := strings.SplitN(body, "\n", 3) // Split into max 3 lines

    if len(lines) < 2 {
        return "", time.Time{}, "", fmt.Errorf("invalid file format: missing metadata lines")
    }

    author := strings.TrimSpace(lines[0]) // Trim leading/trailing whitespace
    author = author[1:len(author)-1]      // Remove leading/trailing brackets

    dateCreated := strings.TrimSpace(lines[1])
    dateCreated = dateCreated[1:len(dateCreated)-1]

    // Parse date into time.Time (assuming specific format)
    t, err := parseDate(dateCreated)
    if err != nil {
        return author, time.Time{}, "", fmt.Errorf("invalid date format: %s, error: %w", dateCreated, err)
    }

    // Extract displayed data (assuming it starts from line 3)
    displayedData := ""
    if len(lines) > 2 {
        displayedData = strings.Join(lines[2:], "\n") // Join remaining lines
    }

    return author, t, displayedData, nil
}

func parseDate(dateStr string) (time.Time, error) {
    // Replace the placeholder validation with your actual date format
    // Here's an example for YYYY-MM-DD format:
    layout := "2006-01-02" // Adjust the layout based on your actual date format
    return time.Parse(layout, dateStr)
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
