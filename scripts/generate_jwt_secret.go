package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	// Generate 32 random bytes (256 bits)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating random bytes: %v\n", err)
		os.Exit(1)
	}

	// Encode to base64
	secret := base64.URLEncoding.EncodeToString(bytes)

	fmt.Println("Generated JWT_SECRET:")
	fmt.Println(secret)
	fmt.Println("\nAdd this to your .env file:")
	fmt.Printf("JWT_SECRET=%s\n", secret)
	fmt.Println("\nOr set it as an environment variable:")
	fmt.Printf("export JWT_SECRET=%s\n", secret)
}

