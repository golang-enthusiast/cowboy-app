package domain

import "context"

// JSON field names.
const (
	JSONFieldName   = "name"
	JSONFieldHealth = "health"
	JSONFieldDamage = "damage"
)

// Cowboy - represents a cowboy struct.
type Cowboy struct {
	Name   string `json:"name"`
	Health int32  `json:"health"`
	Damage int32  `json:"damage"`
}

// CowboyRepository - provides access to a storage.
type CowboyRepository interface {
	FindByName(name string) (*Cowboy, error)
	List(searchCriteria *CowboySearchCriteria) ([]*Cowboy, error)
	UpdateHealthPoints(name string, health int32) error
}

// CowboySearchCriteria - cowboys search criteria.
type CowboySearchCriteria struct {
	ExcludeName string
	HealthGt    int32
	Limit       int64
}

// CowboyService - provides access to a business logic.
type CowboyService interface {

	// GetRandomTarget - selects random target.
	GetRandomTarget() (*Cowboy, error)

	// PrepareGunsAndShoot - each cowboys selects random target and shoots.
	PrepareGunsAndShoot(ctx context.Context) error

	// CommitShooting - commit's cowboy shooting.
	CommitShooting(shooterName string, damage int32) (int32, error)
}

// CowboyRandomGenerator - random number generator.
type CowboyRandomGenerator interface {
	GetRandom(n int) int
}
