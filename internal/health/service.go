package health

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthStatus struct {
	Status    string    `json:"status"`
	DB        string    `json:"db"`
	Uptime    string    `json:"uptime"`
	StartedAt time.Time `json:"started_at"`
}

type Service interface {
	Check(ctx context.Context) (HealthStatus, bool)
	StartMonitor(ctx context.Context, interval time.Duration)
}

type service struct {
	pool      *pgxpool.Pool
	startedAt time.Time
	dbUp      atomic.Bool
}

func NewService(pool *pgxpool.Pool) Service {
	s := &service{
		pool:      pool,
		startedAt: time.Now(),
	}
	s.dbUp.Store(pool.Ping(context.Background()) == nil)
	return s
}

// StartMonitor runs a background goroutine that pings the DB every interval.
// It logs a warning when connectivity is lost and an info message when it is
// restored. The goroutine exits when ctx is cancelled.
func (s *service) StartMonitor(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := s.pool.Ping(ctx)
				isUp := err == nil
				wasUp := s.dbUp.Swap(isUp)

				switch {
				case wasUp && !isUp:
					slog.Warn("Database connection lost", "error", err)
				case !wasUp && isUp:
					slog.Info("Database connection restored")
				}
			}
		}
	}()
}

func (s *service) Check(ctx context.Context) (HealthStatus, bool) {
	err := s.pool.Ping(ctx)
	healthy := err == nil
	s.dbUp.Store(healthy)

	dbStatus := "up"
	status := "healthy"
	if !healthy {
		dbStatus = "down"
		status = "degraded"
	}

	return HealthStatus{
		Status:    status,
		DB:        dbStatus,
		Uptime:    time.Since(s.startedAt).Round(time.Second).String(),
		StartedAt: s.startedAt,
	}, healthy
}
