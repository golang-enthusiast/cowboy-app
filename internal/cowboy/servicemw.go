package cowboy

import (
	"context"
	"cowboy-app/internal/domain"

	"github.com/go-kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(domain.CowboyService) domain.CowboyService

func loggingServiceMiddleware(logger log.Logger) Middleware {
	return func(next domain.CowboyService) domain.CowboyService {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   domain.CowboyService
}

func (mw loggingMiddleware) GetRandomTarget() (target *domain.Cowboy, err error) {
	defer func() {
		_ = mw.logger.Log("method", "GetRandomTarget",
			"target", target,
			"err", err)
	}()
	return mw.next.GetRandomTarget()
}

func (mw loggingMiddleware) PrepareGunsAndShoot(ctx context.Context) (err error) {
	defer func() {
		_ = mw.logger.Log("method", "PrepareGunsAndShoot",
			"err", err)
	}()
	return mw.next.PrepareGunsAndShoot(ctx)
}

func (mw loggingMiddleware) CommitShooting(shooterName string, damage int32) (err error) {
	defer func() {
		_ = mw.logger.Log("method", "CommitShooting",
			"shooterName", shooterName,
			"damage", damage,
			"err", err)
	}()
	return mw.next.CommitShooting(shooterName, damage)
}
