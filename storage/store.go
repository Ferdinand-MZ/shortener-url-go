package storage

import (
	"context"
	"database/sql"
	"os"
	"time"
    "log"
    _ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

// Store represents an in-memory storage for URL mappings
// tapi diganti sekarang pakai postgresql
type Store struct {
    // urls map[string]string // Key : Short URL, Value: Long URL || gapake in memory lagi
    db  *sql.DB
    // mu  sync.RWMutex // Mutex for concurrent access: 1 writer and multiple readers || tapi ga dipake karena pakai postgresql
    redisClient *redis.Client    
}

// NewStore creates a new in-memory store
func NewStore() *Store{

    db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
    
    if err != nil {
        log.Fatal(err)
    }

    return &Store{
        db: db,
        redisClient: redis.NewClient(&redis.Options{
            Addr: os.Getenv("REDIS_ADDR"),
            Password: "",
            DB: 0,
            Protocol: 2,
        }),
    }
}

//Save adds a new URL mapping to the store 
func (s *Store) Save(shortURL, longURL string) error{
    ctx := context.Background()

    // log.Printf("Saving shortURL: %v", shortURL) // buat debugging

    // udah ga pakai mu
    // s.mu.Lock() // Lock for writing 
    // defer s.mu.Unlock() // Unlock after writing

    err := s.redisClient.Set(ctx, shortURL, longURL, 24*time.Hour).Err()
    if err != nil {
        panic(err)
    }

    // s.urls[shortURL] = longURL // Save the mapping 
    // gapake in memory lagi, pake sql:
    _, err2 := s.db.ExecContext(ctx, "INSERT INTO urls (short_url, long_url) VALUES ($1, $2)", shortURL, longURL)
    if err2 != nil {
        // log.Fatal(err2) // yang bikin error
        return err2
    }
    return nil // kalau ga ada error
}

// Get Retreives a long URL for a given short URL 
func (s *Store) Get(shortURL string)(string, bool){
    ctx := context.Background()

    // gapake in memory
    // s.mu.RLock() // Lock for reading 
    // defer s.mu.RUnlock() // Unlock after reading
    
    val, err := s.redisClient.Get(ctx, shortURL).Result()    
    // cek redis
    if err == nil {
        // kalau miss, return langsung
        return val, true
    }

    // kalau hit, buat ini
    // longURL, exists := // s.db[shortURL] 
    // ExecContext untuk insert/update/delete

    var longURL string
    err2 :=  s.db.QueryRowContext(ctx, "SELECT long_url FROM urls WHERE short_url = $1", shortURL).Scan(&longURL)
    if err2 == sql.ErrNoRows {
        log.Printf("Ga ada user dengan yang terdeteksi")
        return "", false // wajib tetep di return biar ga langsung ke kondisi bawah
    } else if err2 != nil {
        log.Fatalf("query nya error : %v\n", err2)
    }

    return longURL, true
}

// Delete  
func (s *Store) Delete(shortURL string){
    ctx := context.Background()

    // redis ini
    val, err := s.redisClient.Del(ctx, shortURL).Result()

    // log.Printf("Redis Del result: %v, err: %v", val, err) // buat debugging 
    
    if err != nil {
        log.Printf("%v", val)
    }

    // var longURL string
    _, err2 :=  s.db.ExecContext(ctx, "DELETE FROM urls WHERE short_url = $1", shortURL)
        if err2 != nil {
        log.Fatal(err2)
    }
}