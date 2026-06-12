package storage

import (
	"testing"
    "github.com/joho/godotenv"
)

func TestSave(t *testing.T) {

	// load env
	godotenv.Load("../.env")

	// arrange
	store := NewStore()
	shortURL := "abc123"
	longURL := "https://google.com"

	// Act
	store.Save(shortURL, longURL)
	hasil, ada := store.Get(shortURL)

	// Assert
	if ada == false || hasil != longURL {
		t.Errorf("Error: got %v, want %v", hasil, longURL)
	}

	t.Cleanup(func() {
		store.Delete(shortURL)

	})
}

func TestGet(t *testing.T) {

	// load env
	godotenv.Load("../.env")

	// arrange
	store := NewStore()
	shortURL := "abc123"
	longURL := "https://google.com" // untuk dibandingkan di assert

	// Act
	store.Save(shortURL, longURL)
	hasil, ada := store.Get(shortURL)

	// Assert
	if ada == false || hasil != longURL {
		t.Errorf("Error: got %v, want %v", hasil, longURL)
	}

	t.Cleanup(func() {
		store.Delete(shortURL)
	})
}