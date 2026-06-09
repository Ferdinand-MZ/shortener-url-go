package models 

import (
 "crypto/rand"
 "encoding/base64"
)

// GenerateShortURL creates a unique 6-character short URL 
func GenerateShortURL() string {
 randomBytes := make([]byte, 6) // Create a byte slice(array) of length 6
 _, err := rand.Read(randomBytes) // Fill the byte slice with random bytes

 if err != nil {
  panic(err) // Handle error (in production, we might want to return an error instead)
 }

 return base64.URLEncoding.EncodeToString(randomBytes)[:6] // Encode the bytes to a base64 string and return the first 6 characters
}