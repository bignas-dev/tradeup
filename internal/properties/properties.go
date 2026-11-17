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

