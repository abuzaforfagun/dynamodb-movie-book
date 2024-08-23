package core_models

import (
	"errors"
	"strings"
)

type Genre int

const (
	NotSupported Genre = -1 + iota
	Unknown
	Romantic
	Action
	Drama
	Travel
	History
)

func (genre Genre) ToString() string {
	return [...]string{"Unknown", "Romantic", "Action", "Drama", "Travel", "History"}[genre]
}

func ToGenre(s string) (Genre, error) {
	switch strings.ToLower(s) {
	case "romantic":
		return Romantic, nil
	case "action":
		return Action, nil
	case "drama":
		return Drama, nil
	case "travel":
		return Travel, nil
	case "history":
		return History, nil
	case "unknown":
		return Unknown, nil
	default:
		return NotSupported, errors.New("invalid genre")
	}
}
