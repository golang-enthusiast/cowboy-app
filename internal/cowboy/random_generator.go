package cowboy

import (
	"cowboy-app/internal/domain"
	"math/rand"
	"time"
)

type randsvc struct {
}

// NewCowboyRandomGenerator - new random generator.
func NewCowboyRandomGenerator() domain.CowboyRandomGenerator {
	return &randsvc{}
}

func (r *randsvc) GetRandom(n int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(n)
}
