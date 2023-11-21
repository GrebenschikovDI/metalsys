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

// PgStorage represents a PostgreSQL storage implementation.
type PgStorage struct {
	db      *sql.DB
	connStr string
}

// InitDB initializes a new database connection using the provided connection
// string and runs migrations located at the specified directory path.
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

// runMigrations applies database migrations from the given directory path.
func (s *PgStorage) runMigrations(connStr, dirPath string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", dirPath), connStr)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// CreateMetricsTable creates the 'metrics' table in the database if it doesn't exist.
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

// AddMetric inserts a single metric into the database.
// If a metric with the same ID already exists, it updates the existing record.
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

// AddMetrics performs a batch insert of multiple metrics into the database.
// Similar to AddMetric, it updates the metrics if they already exist.
// The operation is performed within a transaction.
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

// GetMetric retrieves a single metric by its ID from the database.
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

// GetMetrics fetches all metrics from the database.
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
