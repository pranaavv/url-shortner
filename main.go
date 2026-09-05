package main

import (
	"fmt"
	"time"
)

type URL struct {
	ID          string
	OriginalURL string
	ShortURL    string
	CreatedAt   time.Time
}

func main() {
	fmt.Println("URL Shortner Starting...")
}
