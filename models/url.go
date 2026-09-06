package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	// ErrURLNotFound is returned when a requested short URL does not exist.
	ErrURLNotFound = errors.New("url not found")
	// ErrEmptyURL is returned when the provided original URL is empty.
	ErrEmptyURL = errors.New("url cannot be empty")
	// ErrInvalidURL is returned when the provided original URL is malformed.
	ErrInvalidURL = errors.New("invalid url format, must include http:// or https://")
)

// URL represents a shortened URL record.
type URL struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// URLStore manages in-memory storage of URLs in a concurrency-safe manner.
type URLStore struct {
	mu sync.RWMutex
	db map[string]URL
}

// NewURLStore initializes and returns a new URLStore.
func NewURLStore() *URLStore {
	return &URLStore{
		db: make(map[string]URL),
	}
}

// GenerateShortID computes an 8-character MD5 hash of the original URL.
func GenerateShortID(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:8]
}

// Create validates the input URL, generates a short ID, and saves the URL record.
func (s *URLStore) Create(originalURL string) (*URL, error) {
	trimmed := strings.TrimSpace(originalURL)
	if trimmed == "" {
		return nil, ErrEmptyURL
	}

	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return nil, ErrInvalidURL
	}

	id := GenerateShortID(trimmed)

	record := URL{
		ID:          id,
		OriginalURL: trimmed,
		ShortURL:    id,
		CreatedAt:   time.Now(),
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.db[id] = record

	return &record, nil
}

// GetByID retrieves a URL record by its short ID.
func (s *URLStore) GetByID(id string) (*URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, exists := s.db[id]
	if !exists {
		return nil, ErrURLNotFound
	}

	return &record, nil
}
