// Package tradeups provides logic for generating best CS2 tradeups
package tradeups

import (
	"tradeup/internal/prices"
	"tradeup/internal/properties"
)

// Calculator computes tradeup prices and outcomes from skin prices and collection data.
type Calculator struct {
	priceMap               map[prices.PriceKey]prices.PriceType
	rarityCollectionLookup map[properties.RarityNumberType]map[properties.Collection][]properties.Property
}

const (
	epsilon  = 0.0001
	infinity = 1_000_000_000
)

type inputInfo struct {
	name             properties.SkinName
	collection       properties.Collection
	wear             properties.WearType
	minRange         properties.FloatType
	maxRange         properties.FloatType
	price            prices.PriceType
	minFloat         properties.FloatType
	maxFloat         properties.FloatType
	adjustedMinFloat properties.FloatType
	adjustedMaxFloat properties.FloatType
	skinType         properties.SkinType
}

type outputInfo struct {
	float         properties.FloatType
	avgPriceTaxed prices.PriceType
	taxedPrices   []prices.PriceType
}

type unionColRar struct {
	collection properties.Collection
	rarityNum  properties.RarityNumberType
}

type tradeupOption struct {
	collection properties.Collection
	inputs     []inputInfo
	outputs    []outputInfo
}

// Input is a single user-selected item contributing to a tradeup.
type Input struct {
	Amount        int
	SkinName      string
	Collection    string
	Wear          string
	Float         float64
	AdjustedFloat float64
	Price         int
	SkinType      string
	Rarity        string
}

// Output is one possible outcome of a tradeup with its odds and economics.
type Output struct {
	Wear        string
	Collection  string
	SkinName    string
	Float       float64
	MaxFloat    float64
	MinFloat    float64
	Price       int
	PriceTaxed  int
	ProfitTaxed float64
	Chance      float64
	SkinType    string
	Profit      int
}

// Tradeup is a persisted tradeup candidate with its inputs and outcomes.
type Tradeup struct {
	ID            int
	Profitability float64
	FreeFloat     float64
	RarityNum     int
	Rarity        string
	TradeupType   string
	AvgFloat      float64
	Price         int
	Return        int
	NetProfit     int
	ProfitChance  float64
	UniqueInputs  int

	Input1Amount        int
	Input1Collection    string
	Input1Skin          string
	Input1Wear          string
	Input1Float         float64
	Input1AdjustedFloat float64
	Input1Price         int
	Input1SkinType      string

	Input2Amount        int
	Input2Collection    string
	Input2Skin          string
	Input2Wear          string
	Input2Float         float64
	Input2AdjustedFloat float64
	Input2Price         int
	Input2SkinType      string

	Outcomes []Output
}
