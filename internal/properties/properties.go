// Package properties models skin, wear, rarity and collection metadata, and syncs it to the database.
package properties

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type (
	SkinName string
	SkinType string
)

const (
	InvalidType  = SkinType("invalid type")
	NormalType   = SkinType("Normal")
	SouvenirType = SkinType("Souvenir")
	StatTrakType = SkinType("StatTrak™")
	MixedType    = SkinType("Mixed")
)

const (
	StatTrakPrefix = string(StatTrakType) + " "
	SouvenirPrefix = string(SouvenirType) + " "
)

// FloatType CS2 stores floats as f32
type FloatType float32

type WearType string

type RarityNumberType int

type RarityStringType string

const (
	InvalidRarity   = RarityStringType("error rarity")
	ConsumerGrade   = RarityStringType("Consumer Grade")
	IndustrialGrade = RarityStringType("Industrial Grade")
	MilSpecGrade    = RarityStringType("Mil-Spec Grade")
	Restricted      = RarityStringType("Restricted")
	Classified      = RarityStringType("Classified")
	Covert          = RarityStringType("Covert")
	Extraordinary   = RarityStringType("Extraordinary")
)

var Rarities = []RarityStringType{
	InvalidRarity,
	ConsumerGrade,
	IndustrialGrade,
	MilSpecGrade,
	Restricted,
	Classified,
	Covert,
	Extraordinary,
}

type WearData struct {
	Name     WearType
	MinFloat FloatType
	MaxFloat FloatType
}

const InvalidWear = WearType("error wear")

var Wears = []WearData{
	{InvalidWear, -9999, -999},
	{"Factory New", 0.00, 0.07},
	{"Minimal Wear", 0.07, 0.15},
	{"Field-Tested", 0.15, 0.38},
	{"Well-Worn", 0.38, 0.45},
	{"Battle-Scarred", 0.45, 1.00},
}

func ExtractWearFromMarketName(n string) (WearType, error) {
	for _, w := range Wears {
		if strings.Contains(n, string(w.Name)) {
			return w.Name, nil
		}
	}
	return "", nil
}

type Map map[SkinName]Property

var (
	dopplerPhases      = []string{"Phase 1", "Phase 2", "Phase 3", "Phase 4", "Black Pearl", "Saphire", "Ruby"}
	gammaDopplerPhases = []string{"Phase 1", "Phase 2", "Phase 3", "Phase 4", "Emerald"}
)

func ExtractRarityNumberFromRarityString(rs string) (int, error) {
	switch rs {
	case "Consumer Grade":
		return 1, nil
	case "Industrial Grade":
		return 2, nil
	case "Mil-Spec Grade":
		return 3, nil
	case "Restricted":
		return 4, nil
	case "Classified":
		return 5, nil
	case "Covert":
		return 6, nil
	case "Extraordinary":
		return 7, nil
	default:
		return 0, fmt.Errorf("invalid rarity string %q", rs)
	}
}

// FilterValue is a single skin-filter suggestion result.
type FilterValue struct {
	Name       string `json:"name"`
	Collection string `json:"collection"`
}

const skinsURL = "https://raw.githubusercontent.com/ByMykel/CSGO-API/main/public/api/en/skins.json"

var replacer = strings.NewReplacer(
	"The ", "",
	" Collection", "",
	" Case", "",
	"Operation ", "",
	" Weapon", "",
)

type rarity struct {
	Name string `json:"name"`
}

type origin struct {
	Name string `json:"name"`
}

type skinResponse struct {
	Name        string   `json:"name"`
	MinFloat    float64  `json:"min_float"`
	MaxFloat    float64  `json:"max_float"`
	Rarity      rarity   `json:"rarity"`
	Stattrak    bool     `json:"stattrak"`
	Collections []origin `json:"collections"`
	Crates      []origin `json:"crates"`
}

type Collection string

type Property struct {
	Name        SkinName
	Collection  Collection
	Rarity      string
	RarityNum   int
	Tradeupable bool
	Stattrak    bool
	MinFloat    float64
	MaxFloat    float64
}

// GetProperties fetches the skin list and upserts each skin's properties into the database.
func GetProperties() (Map, error) {
	resp, err := http.Get(skinsURL)
	if err != nil {
		return nil, fmt.Errorf("http.Get failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non ok status code: %d", resp.StatusCode)
	}

	var skins []skinResponse
	if err := json.NewDecoder(resp.Body).Decode(&skins); err != nil {
		return nil, fmt.Errorf("json decoding failed: %w", err)
	}

	properties := make(Map)

	for _, skin := range skins {
		if skin.Name == "" {
			continue
		}

		if strings.HasPrefix(skin.Name, "★") {
			skin.Rarity.Name = "Extraordinary"
		}

		collections := getCollections(skin.Collections, skin.Crates)
		phases := []string{""}

		if strings.Contains(skin.Name, "Gamma Doppler") {
			phases = gammaDopplerPhases
		} else if strings.Contains(skin.Name, "Doppler") {
			phases = dopplerPhases
		}

		for _, phase := range phases {
			for _, col := range collections {
				name := skin.Name

				if phase != "" {
					name = fmt.Sprintf("%s %s", name, phase)
				}

				rarNum, err := ExtractRarityNumberFromRarityString(skin.Rarity.Name)
				if err != nil {
					return nil, fmt.Errorf("error some skin doesnt have rarity num: %w", err)
				}

				properties[SkinName(name)] = Property{
					Name:        SkinName(name),
					Collection:  Collection(col),
					Rarity:      skin.Rarity.Name,
					MinFloat:    skin.MinFloat,
					MaxFloat:    skin.MinFloat,
					RarityNum:   rarNum,
					Tradeupable: false,
					Stattrak:    skin.Stattrak,
				}
			}
		}
	}
	return properties, nil
}

func getCollections(args ...[]origin) []string {
	seen := make(map[string]struct{})
	for _, group := range args {
		for _, col := range group {
			cleanCol := replacer.Replace(col.Name)
			if cleanCol == "Limited Edition Item" {
				continue
			}
			seen[cleanCol] = struct{}{}
		}
		if len(seen) != 0 {
			break
		}
	}

	result := make([]string, 0, len(seen))
	for k := range seen {
		result = append(result, k)
	}

	return result
}

// GetWearName returns the wear-condition name for a float value.
func GetWearName(f float64) string {
	if f < 0 {
		return string(Wears[1].Name)
	}
	for i := 1; i < len(Wears); i++ {
		if f < float64(Wears[i].MaxFloat) {
			return string(Wears[i].Name)
		}
	}
	return string(Wears[len(Wears)-1].Name)
}
