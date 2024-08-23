package core_models

import "strings"

var supportedGenre = map[string]bool{
	"action":      true,
	"adventure":   true,
	"animation":   true,
	"biography":   true,
	"comedy":      true,
	"crime":       true,
	"documentary": true,
	"drama":       true,
	"family":      true,
	"fantasy":     true,
	"film-noir":   true,
	"history":     true,
	"horror":      true,
	"musical":     true,
	"mystery":     true,
	"romance":     true,
	"sci-fi":      true,
	"sport":       true,
	"thriller":    true,
	"war":         true,
	"western":     true,
}

func IsSupportedGenre(val string) bool {
	return supportedGenre[strings.ToLower(val)]
}
