// Package prices fetches Steam market prices and provides price/economy helpers.
package prices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"tradeup/internal/properties"
)

const latestPricesURL = "https://raw.githubusercontent.com/ByMykel/counter-strike-price-tracker/refs/heads/main/static/latest.json"

func ApplySteamFee(p PriceType) PriceType {
	return p * 20 / 23
}

type latestPrices struct {
	Prices map[string]int64 `json:"prices"`
}

type PriceType int

type PriceKey struct {
	Name     properties.SkinName
	Wear     properties.WearType
	SkinType properties.SkinType
}

// GetPrices gets all of the prices
func GetPrices(skinProperties properties.Map) (map[PriceKey]PriceType, error) {
	resp, err := http.Get(latestPricesURL)
	if err != nil {
		return nil, fmt.Errorf("fetching latest prices: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching latest prices: HTTP %d", resp.StatusCode)
	}

	var data latestPrices
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decoding latest prices: %w", err)
	}
	if len(data.Prices) == 0 {
		return nil, fmt.Errorf("latest prices payload is empty")
	}

	priceLookup := make(map[PriceKey]PriceType, len(data.Prices))
	sortedSkins := make([]properties.SkinName, 0, len(skinProperties))
	for key := range data.Prices {
		sortedSkins = append(sortedSkins, properties.SkinName(key))
	}
	sort.Slice(sortedSkins, func(i, j int) bool {
		return len(sortedSkins[i]) > len(sortedSkins[j])
	})

	for n, p := range data.Prices {
		rawName := n
		cleanName, souvenir := strings.CutPrefix(rawName, properties.SouvenirPrefix)
		cleanName, stattrak := strings.CutPrefix(cleanName, properties.StatTrakPrefix)

		if souvenir && stattrak {
			return nil, fmt.Errorf("skin is statttrak and souvenir at the same time: %v", rawName)
		}

		skinType := properties.NormalType
		if souvenir {
			skinType = properties.SouvenirType
		} else if stattrak {
			skinType = properties.StatTrakType
		}

		skinName := properties.SkinName("")
		for _, skin := range sortedSkins {
			if strings.HasPrefix(cleanName, string(skin)) {
				skinName = skin
				break
			}
		}
		if skinName == "" {
			continue
		}

		wear, err := properties.ExtractWearFromMarketName(cleanName)
		if err != nil {
			return nil, fmt.Errorf("extracting wear from %q: %w", cleanName, err)
		}

		priceLookup[PriceKey{
			Name:     skinName,
			Wear:     wear,
			SkinType: skinType,
		}] = PriceType(p)
	}

	return priceLookup, nil
}
