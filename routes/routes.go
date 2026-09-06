package routes

import (
	"net/http"

	"url-shortner/controllers"
)

// SetupRoutes registers application routes and maps them to their respective controller handlers.
func SetupRoutes(ctrl *controllers.URLController) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", ctrl.Root)
	mux.HandleFunc("/shorten", ctrl.Shorten)
	mux.HandleFunc("/blob/", ctrl.Redirect)

	return mux
}
