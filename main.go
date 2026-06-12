
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

    ctx := context.Background()

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
    // port := ":8080"
    println("Server is running on port", s.Addr)

    // buat channel untuk jadi medium pengirim sinyal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)


    // go s.ListenAndServe()

    // if err := http.ListenAndServe(port, rateHandler); err != nil {
    // ini anonymouse function makanya ga punya nama
    // ini udah dibuat goroutine juga
    go func() {
        if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        println("Error starting server:", err)
        } else if err == http.ErrServerClosed{
            println("Server berhasil dimatikan")
        }
    } ()

    <- quit

    if err := s.Shutdown(ctx); err != nil {
        println("Error shutdown, ada kejanggalan: ", err)
    }

}