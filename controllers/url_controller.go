package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"url-shortner/models"
	"url-shortner/views"
)

// URLController handles incoming HTTP requests for URL operations.
type URLController struct {
	store *models.URLStore
}

// NewURLController creates a new URLController instance with the given store.
func NewURLController(store *models.URLStore) *URLController {
	return &URLController{
		store: store,
	}
}

// Root handles requests to the root path ("/").
func (c *URLController) Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		views.RenderError(w, http.StatusNotFound, "route not found")
		return
	}

	if r.Method != http.MethodGet {
		views.RenderError(w, http.StatusMethodNotAllowed, "method not allowed, GET required")
		return
	}

	views.RenderJSON(w, http.StatusOK, views.MessageResponse{
		Message: "hello from server",
	})
}

// Shorten handles the creation of a shortened URL from a JSON payload.
func (c *URLController) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		views.RenderError(w, http.StatusMethodNotAllowed, "method not allowed, POST required")
		return
	}

	var req struct {
		URL string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		views.RenderError(w, http.StatusBadRequest, "invalid request body, valid JSON with 'url' required")
		return
	}

	record, err := c.store.Create(req.URL)
	if err != nil {
		if errors.Is(err, models.ErrEmptyURL) || errors.Is(err, models.ErrInvalidURL) {
			views.RenderError(w, http.StatusBadRequest, err.Error())
			return
		}
		views.RenderError(w, http.StatusInternalServerError, "failed to create short url")
		return
	}

	views.RenderJSON(w, http.StatusCreated, views.ShortenResponse{
		ID:          record.ID,
		OriginalURL: record.OriginalURL,
		ShortURL:    record.ShortURL,
		CreatedAt:   record.CreatedAt,
	})
}

// Redirect handles redirecting short URLs to their original destination.
func (c *URLController) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		views.RenderError(w, http.StatusMethodNotAllowed, "method not allowed, GET required")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/blob/")
	id = strings.TrimSpace(id)

	if id == "" {
		views.RenderError(w, http.StatusBadRequest, "short url id is required")
		return
	}

	record, err := c.store.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrURLNotFound) {
			views.RenderError(w, http.StatusNotFound, "short url not found")
			return
		}
		views.RenderError(w, http.StatusInternalServerError, "failed to retrieve url")
		return
	}

	http.Redirect(w, r, record.OriginalURL, http.StatusFound)
}
