package storages

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

type PgStorage struct {
	db      *sql.DB
	connStr string
}

func InitDB(_ context.Context, connStr, dirPath string) (*PgStorage, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	storage := &PgStorage{
		db:      db,
		connStr: connStr,
	}
	err = storage.runMigrations(connStr, dirPath)
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func (s *PgStorage) runMigrations(connStr, dirPath string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", dirPath), connStr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func (s *PgStorage) CreateMetricsTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS metrics (
    		id TEXT PRIMARY KEY,
    		mtype TEXT,
    		delta BIGINT,
    		value DOUBLE PRECISION
		);
	`)
	return err
}

func (s *PgStorage) AddMetric(ctx context.Context, mc models.Metric) error {
	query := `
		INSERT INTO metrics (id, mtype, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET delta = metrics.delta + EXCLUDED.delta, value = EXCLUDED.value;
	`
	_, err := s.db.ExecContext(ctx, query, mc.ID, mc.Mtype, mc.Delta, mc.Value)
	return err
}

func (s *PgStorage) AddMetrics(ctx context.Context, metrics []models.Metric) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO metrics (id, mtype, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET delta = metrics.delta + EXCLUDED.delta, value = EXCLUDED.value;
	`

	for _, mc := range metrics {
		_, err := tx.ExecContext(ctx, query, mc.ID, mc.Mtype, mc.Delta, mc.Value)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *PgStorage) GetMetric(ctx context.Context, name string) (value models.Metric, err error) {
	row := s.db.QueryRowContext(ctx, "SELECT id, mtype, delta, value FROM metrics WHERE id = $1", name)

	err = row.Scan(&value.ID, &value.Mtype, &value.Delta, &value.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = pgx.ErrNoRows
		}
	}
	return value, err
}

func (s *PgStorage) GetMetrics(ctx context.Context) ([]models.Metric, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, mtype, delta, value FROM metrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.Metric
	for rows.Next() {
		var metric models.Metric
		if err := rows.Scan(&metric.ID, &metric.Mtype, &metric.Delta, &metric.Value); err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return metrics, nil
}
