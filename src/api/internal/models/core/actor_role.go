package core_models

import (
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/custom_errors"
)

type ActorRole int

const (
	LeadHero ActorRole = iota
	LeadHeroin
	LeadBillen
	Hero
	Heroin
	Billen
	Other
)

func (role ActorRole) ToString() (string, error) {
	if role < LeadHero || role > Other {
		return "", &custom_errors.BadRequestError{
			Message: "Please provide a valid actor role",
		}
	}
	return [...]string{"Lead Hero", "Lead Heroin", "Lead Billen", "Hero", "Heroin", "Billen", "Other"}[role], nil
}
