package main

import (
	"bruvela-backend/pkg/auth"
	"fmt"
)

func main() {
	hash, err := auth.HashPassword("admin123")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(hash)
}
