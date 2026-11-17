package main

import (
	"fmt"
	"log"

	"tradeup/internal/currencies"
	"tradeup/internal/images"
)

func main() {
	fmt.Println("Hello, World!")
	if err := images.SaveImages(); err != nil {
		log.Fatal(err)
	}

	data, err := currencies.GetCurrencies()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
}
