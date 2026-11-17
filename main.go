package main

import (
	"fmt"
	"log"

	"tradeup/internal/currencies"
	"tradeup/internal/images"
	"tradeup/internal/prices"
	"tradeup/internal/properties"
)

func main() {
	if err := images.SaveImages(); err != nil {
		log.Fatal(err)
	}

	properties, err := properties.GetProperties()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", properties)

	currencies, err := currencies.GetCurrencies()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", currencies)

	prices, err := prices.GetPrices(properties)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", prices)
}
