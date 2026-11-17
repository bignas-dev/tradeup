// Package nameformatter is for systemizing property fields
package nameformatter

import "strings"

var collectionReplacer = strings.NewReplacer(
	"The ", "",
	" Collection", "",
	" Case", "",
	"Operation ", "",
	" Weapon", "",
)

func FormatCollectionString(collection string) string {
	if collection == "Limited Edition Item" {
		return ""
	}
	return collectionReplacer.Replace(collection)
}
