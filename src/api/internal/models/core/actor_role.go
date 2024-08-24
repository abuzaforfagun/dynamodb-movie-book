package core_models

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

func (status ActorRole) ToString() string {
	return [...]string{"Lead Hero", "Lead Heroin", "Lead Billen", "Hero", "Heroin", "Billen", "Other"}[status]
}
