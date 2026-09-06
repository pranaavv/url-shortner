package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	fmt.Println("Hasher:- ", hasher)
	data := hasher.Sum(nil)
	fmt.Println("Data: ", data)
	hash := hex.EncodeToString(data)
	fmt.Println("Hash: ", hash)
	fmt.Println("Hash: ", hash[:8])

	return hash[:8]
}

func createURL(OriginalURL string) string {
	shortURL := generateShortURL(OriginalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:          id,
		OriginalURL: OriginalURL,
		ShortURL:    shortURL,
		CreatedAt:   time.Now(),
	}

	return shortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("Error n finding URL, func :- getURL")
	}
	return url, nil
}

func handler(w http.ResponseWriter,r *http.Request){
	fmt.Fprintln(w,"hello from server")
}




func main() {
	fmt.Println("URL Shortner Starting...")
	OriginalURL := "https://linkedin.com/in/"
	fmt.Println("URL: ")
	generateShortURL(OriginalURL)
	// fmt.Println("URL: ",ShortURL)

	http.HandleFunc("/",handler)

	fmt.Println("Starting the server ....")
	err := http.ListenAndServe(":3000",nil)
	if err!=nil{
		fmt.Println("Error while starting server :- ",err)
	}
}
