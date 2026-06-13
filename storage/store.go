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

type Store struct {
    db  *sql.DB
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

    err := s.redisClient.Set(ctx, shortURL, longURL, 24*time.Hour).Err()
    // err := s.redisClient.Set(ctx, shortURL, longURL, 10*time.Second).Err() // buat testing
    log.Printf("shortURL: %v", shortURL)
    if err != nil {
        panic(err)
    }



    // _, err2 := s.db.ExecContext(ctx, "INSERT INTO urls (short_url, long_url, expires_at) VALUES ($1, $2, NOW() + INTERVAL '30 seconds')", shortURL, longURL) // for testing
    _, err2 := s.db.ExecContext(ctx, "INSERT INTO urls (short_url, long_url, expires_at) VALUES ($1, $2, NOW() + INTERVAL '24 hour')", shortURL, longURL)
    if err2 != nil {
        return err2
    }
    return nil
}

// Get Retreives a long URL for a given short URL 
func (s *Store) Get(shortURL string)(string, bool){
    ctx := context.Background()


    
    val, err := s.redisClient.Get(ctx, shortURL).Result()   
    ttl, _ := s.redisClient.TTL(ctx, shortURL).Result() // kita ga perlu handling error redis, karena sudah ada di atas variabel err
    log.Printf("Redis err: %v, TTL: %v", err, ttl)
    

    // cek redis
    if err == nil {
        // kalau miss, return langsung
        if ttl <= 0 {
            return "", false
        }
            // ga perlu tambahin ttl > 0
            return val, true
    }

    var longURL string
    err2 :=  s.db.QueryRowContext(ctx, "SELECT long_url FROM urls WHERE short_url = $1 AND expires_at > NOW()", shortURL).Scan(&longURL)
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