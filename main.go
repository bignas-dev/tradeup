package main

import (
	"fmt"
	"log"

	"tradeup/internal/images"
)

func main() {
	fmt.Println("Hello, World!")
	if err := images.SaveImages(); err != nil {
		log.Fatal(err)
	}
}
