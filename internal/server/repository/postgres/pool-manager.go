package postgres

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Settings struct {
	ConnString string
	RetryCount int
	RetryDelay time.Duration
}

type PoolManager struct {
	pool        *pgxpool.Pool
	config      *pgxpool.Config
	retryCount  int
	retryDelay  time.Duration
	reconecting atomic.Bool
}

func NewConnectionPool(ctx context.Context, settings Settings) (*PoolManager, error) {
	config, err := pgxpool.ParseConfig(settings.ConnString)
	if err != nil {
		return nil, err
	}

	pool := &PoolManager{
		config:      config,
		retryCount:  settings.RetryCount,
		retryDelay:  settings.RetryDelay,
		reconecting: atomic.Bool{},
	}

	err = pool.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (p *PoolManager) Pool() *pgxpool.Pool {
	return p.pool
}

func (p *PoolManager) Connect(ctx context.Context) error {
	// If already reconecting, just wait
	if p.reconecting.Load() {
		for p.reconecting.Load() {
			time.Sleep(p.retryDelay)
		}
		return nil
	}

	p.reconecting.Store(true)
	defer p.reconecting.Store(false)

	// Close pool if already opened
	if p.pool != nil {
		p.pool.Close()
	}

	var err error
	var pool *pgxpool.Pool

	// Connecting to db
	for i := range p.retryCount {
		slog.Info(
			"Attempting connect to the DB",
			slog.Int("Attempt", i),
			slog.Int("Max attempts", p.retryCount),
		)
		pool, err = pgxpool.NewWithConfig(ctx, p.config)
		if err == nil {
			p.pool = pool
			return nil
		}
		time.Sleep(p.retryDelay)
	}

	slog.Error("Error connecting to DB", slog.Any("Error", err))

	return err
}
