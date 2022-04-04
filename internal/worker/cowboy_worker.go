package worker

import (
	"context"

	"cowboy-app/internal/domain"

	"github.com/go-kit/log"
)

// CowboyWorker struct.
type CowboyWorker struct {
	*Worker
	cowboyService domain.CowboyService
}

// NewCowboyWorker - sets up a new worker.
func NewCowboyWorker(
	props *Props,
	queueAPI domain.QueueAPI,
	cowboyService domain.CowboyService,
	logger log.Logger) *CowboyWorker {
	var (
		worker = new(props, queueAPI, logger)
	)
	return &CowboyWorker{
		Worker:        worker,
		cowboyService: cowboyService,
	}
}

// HandlePrepareGunsMessage - handle message.
func (w *CowboyWorker) HandlePrepareGunsMessage(
	ctx context.Context,
	_ *domain.PrepareGunsMessage,
) error {
	return w.cowboyService.PrepareGunsAndShoot(ctx)
}

// HandleShootMessage - handle message.
func (w *CowboyWorker) HandleShootMessage(
	ctx context.Context,
	msg *domain.ShootMessage,
) error {
	_, err := w.cowboyService.CommitShooting(msg.ShooterName, msg.Damage)
	return err
}

// HandleWinnerMessage - handle message.
func (w *CowboyWorker) HandleWinnerMessage(
	ctx context.Context,
	msg *domain.WinnerMessage,
) error {
	return w.Logger.Log(msg.Message)
}
