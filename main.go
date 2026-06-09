
package main
import (
    "net/http"
    "url-shortener/handlers"
    "url-shortener/storage"
    "url-shortener/middleware"
    "github.com/joho/godotenv"
    "context"
    "os"
    "os/signal"
)

func main() {
    godotenv.Load()


    //initialise the store
    store := storage.NewStore()

    // Initialise the handler
    handler := handlers.NewHandler(store)

    // inisialisasi rate limiting
    mux := http.NewServeMux()

    // Define Routes
    mux.HandleFunc("/shorten", handler.ShortenURL) // POST /shorten
    mux.HandleFunc("/", handler.RedirectURL)     // GET /:shortURL

    // 10 requests/second, burst of 20
    rateLimiter := middleware.NewRateLimitMiddleware(10, 20)

    // wrapping
    rateHandler := rateLimiter.Middleware(mux)

    s := http.Server{
        Addr: ":8080",
        Handler: rateHandler,
    }

    //start the server 
    port := ":8080"
    println("Server is running on port", port)
    
    if err := http.ListenAndServe(port, rateHandler); err != nil {
    println("Error starting server:", err)
    }
    
    s.ListenAndServe()

}