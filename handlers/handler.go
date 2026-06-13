package handlers

import (
 "fmt"
 "net/http"
//these are our internal packages
 "url-shortener/models" 
 "url-shortener/storage"
)

// Handler contains dependencies for handling HTTP requests
type Handler struct {
    store *storage.Store
}

// NewHandler initialises a new Handler
func NewHandler(store *storage.Store) *Handler {
    return &Handler{store: store}
}

// ShortenURL handles POST requests to shorten a URL
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

// parse the long url and request body
    r.ParseForm()
    longURL := r.FormValue("url")
    if longURL == "" {
        http.Error(w, "missing 'url' parameter", http.StatusBadRequest)
        return
    }

    alias := r.FormValue("alias")

    var shortURL string // ini biar percabangan ngerti shortURL yang kita pakai itu gimana
    // generate short url

    // kalo ga pake keyword alias, contoh : curl -X POST -d "url=https://www.google.com" http://localhost:8080/shorten
    if alias == "" {
        shortURL = models.GenerateShortURL()
        // kalo pake, contoh :  curl -X POST -d "url=https://www.google.com&alias=google" http://localhost:8080/shorten
    } else {
        shortURL = alias
    }
    
    h.store.Save(shortURL, longURL) // save mapping di store

    // respond dengan short URL nya
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("http://shorty.url/" + shortURL)) // should be a real domain name + shortURL

}

// RedirectURL handles GET requests to redirect to the original URL
func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {

    if r.Method != http.MethodGet {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    //extract the short url from request url
    shortURL := r.URL.Path[1:]         // remove leading slash
    fmt.Println("shortURL:", shortURL) // Debugging line

    //retrieve the long url from the store
    longURL, exists := h.store.Get(shortURL)
    if !exists {
        http.Error(w, "URL not found", http.StatusNotFound)
        return
    }

    //Redirect to original URL 
    http.Redirect(w, r, longURL, http.StatusFound)

}