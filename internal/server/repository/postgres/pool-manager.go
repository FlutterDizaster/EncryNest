package postgres

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	//nolint:revieve // This is for migrate
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
)

type Settings struct {
	ConnString          string
	RetryCount          int
	RetryDelay          time.Duration
	MigrationsDirectory string
}

type PoolManager struct {
	pool                *pgxpool.Pool
	config              *pgxpool.Config
	connString          string
	retryCount          int
	retryDelay          time.Duration
	reconecting         atomic.Bool
	migrationsDirectory string
}

func NewPoolManager(ctx context.Context, settings Settings) (*PoolManager, error) {
	config, err := pgxpool.ParseConfig(settings.ConnString)
	if err != nil {
		return nil, err
	}

	pool := &PoolManager{
		config:              config,
		retryCount:          settings.RetryCount,
		connString:          settings.ConnString,
		retryDelay:          settings.RetryDelay,
		reconecting:         atomic.Bool{},
		migrationsDirectory: settings.MigrationsDirectory,
	}

	err = pool.migrate(ctx)
	if err != nil {
		slog.Error("Error migrations", slog.Any("err", err))
		return nil, err
	}

	err = pool.Connect(ctx)
	if err != nil {
		slog.Error("Error while connecting to DB", slog.Any("err", err))
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
			slog.Int("Attempt", i+1),
			slog.Int("Max attempts", p.retryCount),
		)
		pool, err = pgxpool.NewWithConfig(ctx, p.config)
		if err == nil {
			p.pool = pool
			slog.Info("Connected to the DB")
			return nil
		}
		time.Sleep(p.retryDelay)
	}

	slog.Error("Error connecting to DB", slog.Any("Error", err))

	return err
}

func (p *PoolManager) migrate(ctx context.Context) error {
	if p.migrationsDirectory == "" {
		slog.Info("Skipping migrations")
		return nil
	}

	slog.Info("Running migrations")

	migrator, err := migrate.New(
		"file://"+p.migrationsDirectory,
		p.connString,
	)
	if err != nil {
		slog.Error("Error while creating migrator", slog.Any("err", err))
		return err
	}
	defer migrator.Close()

	errCh := make(chan error)

	go func() {
		errCh <- migrator.Up()
	}()

	select {
	case err = <-errCh:
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			slog.Error("Error while running migrations", slog.Any("err", err))
			return err
		}
	case <-ctx.Done():
		migrator.GracefulStop <- true
		return ctx.Err()
	}

	slog.Info("Migrations finished")

	return nil
}
