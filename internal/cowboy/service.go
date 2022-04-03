package cowboy

import (
	"context"
	"fmt"
	"strings"

	"cowboy-app/internal/domain"
	errors "cowboy-app/internal/error"

	"github.com/go-kit/log"
)

type service struct {
	cowboyName      string
	repository      domain.CowboyRepository
	queueService    domain.QueueService
	randomGenerator domain.CowboyRandomGenerator
}

// NewService creates a new service with necessary dependencies.
func NewService(cowboyName string,
	repository domain.CowboyRepository,
	queueService domain.QueueService,
	randomGenerator domain.CowboyRandomGenerator,
	logger log.Logger) domain.CowboyService {
	var service domain.CowboyService
	{
		service = newBasicService(cowboyName, repository, queueService, randomGenerator)
		service = loggingServiceMiddleware(logger)(service)
	}
	return service
}

func newBasicService(cowboyName string,
	repository domain.CowboyRepository,
	queueService domain.QueueService,
	randomGenerator domain.CowboyRandomGenerator) domain.CowboyService {
	return &service{
		cowboyName:      cowboyName,
		repository:      repository,
		queueService:    queueService,
		randomGenerator: randomGenerator,
	}
}

func (s *service) GetRandomTarget() (*domain.Cowboy, error) {
	cowboys, err := s.repository.List(&domain.CowboySearchCriteria{
		ExcludeName: s.cowboyName,
		HealthGt:    0,
		Limit:       100,
	})
	if err != nil {
		return nil, err
	}
	var (
		size = len(cowboys)
	)
	if size == 0 {
		return nil, errors.NewErrNotFound("List of cowboys is empty")
	}
	return cowboys[s.randomGenerator.GetRandom(size)], nil
}

func (s *service) PrepareGunsAndShoot(ctx context.Context) error {
	// Select random target.
	target, err := s.GetRandomTarget()
	if err != nil {
		return err
	}

	// Get current cowboy.
	foundCowboy, err := s.repository.FindByName(s.cowboyName)
	if err != nil {
		return err
	}
	if foundCowboy == nil {
		return errors.NewErrNotFound(
			fmt.Sprintf("Cowboy not found by name: %s", s.cowboyName),
		)
	}

	// Shoot.
	_, err = s.queueService.SendMessage(ctx, target.Name, &domain.ShootMessage{
		ShooterName: foundCowboy.Name,
		Damage:      foundCowboy.Damage,
	})
	return err
}

func (s *service) CommitShooting(shooterName string, damage int32) error {
	// Cowboys don’t shoot themselves.
	if strings.EqualFold(shooterName, s.cowboyName) {
		return errors.NewErrInvalidArgument(
			"Cowboys don’t shoot themselves",
		)
	}

	// Get target cowboy.
	targetCowboy, err := s.repository.FindByName(s.cowboyName)
	if err != nil {
		return err
	}
	if targetCowboy == nil {
		return errors.NewErrNotFound(
			fmt.Sprintf("Cowboy not found by name: %s", s.cowboyName),
		)
	}

	// Cowboys don’t shoot dead cowboys.
	if targetCowboy.Health == 0 {
		return errors.NewErrInvalidArgument(
			"Cowboys don’t shoot dead cowboys.",
		)
	}

	// Calculate health points.
	var remainedHealth int32 = 0
	if targetCowboy.Health >= damage {
		remainedHealth = targetCowboy.Health - damage
	}

	// Update health points.
	return s.repository.UpdateHealthPoints(targetCowboy.Name, remainedHealth)
}
