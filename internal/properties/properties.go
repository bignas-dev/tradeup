// Package properties models skin, wear, rarity and collection metadata, and syncs it to the database.
package properties

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
	InvalidRarity = RarityStringType("error rarity")
	ConsumerGrade = RarityStringType("Consumer Grade")
	IndustrialGrade = RarityStringType("Industrial Grade")
	MilSpecGrade = RarityStringType("Mil-Spec Grade")
	Restricted = RarityStringType("Restricted")
	Classified = RarityStringType("Classified")
	Covert = RarityStringType("Covert")
	Extraordinary = RarityStringType("Extraordinary")
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

