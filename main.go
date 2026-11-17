package main

import (
	"fmt"
	"log"

	"tradeup/internal/currencies"
	"tradeup/internal/images"
	"tradeup/internal/properties"
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
	fmt.Printf("%#v\n", data)

	data2, err := properties.GetProperties()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", data2)
}
