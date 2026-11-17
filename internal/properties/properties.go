// Package properties models skin, wear, rarity and collection metadata, and syncs it to the database.
package properties

import "strings"

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

type Map map[SkinName]Property

var (
	dopplerPhases      = []string{"Phase 1", "Phase 2", "Phase 3", "Phase 4", "Black Pearl", "Saphire", "Ruby"}
	gammaDopplerPhases = []string{"Phase 1", "Phase 2", "Phase 3", "Phase 4", "Emerald"}
)
