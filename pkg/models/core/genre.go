package core_models

type Genre int

const (
	Unknown Genre = iota
	Romantic
	Action
	Drama
	Travel
	History
)

func (genre Genre) ToString() string {
	return [...]string{"Unknown", "Romantic", "Action", "Drama", "Travel", "History"}[genre]
}
