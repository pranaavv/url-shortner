package main

import (
	"fmt"
	"net/http"

	"url-shortner/controllers"
	"url-shortner/models"
	"url-shortner/routes"
)

func main() {
	fmt.Println("URL Shortener starting...")

	store := models.NewURLStore()
	ctrl := controllers.NewURLController(store)
	router := routes.SetupRoutes(ctrl)

	port := ":3000"
	fmt.Printf("Server listening on http://localhost%s\n", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
