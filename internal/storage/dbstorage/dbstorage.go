// Package dbstorage - используется для взаимодействия с БД pgsql
package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"github.com/MaximPolyaev/go-metrics/internal/retry"
)

// timeoutReqs - timeout execute request on DB
const timeoutReqs = 10 * time.Second

type Storage struct {
	db  *sql.DB
	log *logger.Logger
}

func New(db *sql.DB, log *logger.Logger) *Storage {
	return &Storage{db: db, log: log}
}

// Init - init DB storage. Create required table
func (s *Storage) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutReqs)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS metrics (
    id VARCHAR(255) not null,
    type VARCHAR(20) not null,
    delta bigint,
    value double precision,
    CONSTRAINT id_type_unique UNIQUE (id, type)
);
`)

	return err
}

// Set - set metric to DB
func (s *Storage) Set(ctx context.Context, mType metric.Type, val metric.Metric) {
	retry.Retry(s.setHandler(ctx, mType, val), s.log)
}

// BatchSet - batch set metrics to DB
func (s *Storage) BatchSet(ctx context.Context, mSlice []metric.Metric) {
	retry.Retry(s.batchSetHandler(ctx, mSlice), s.log)
}

// Get - get metric from DB
func (s *Storage) Get(ctx context.Context, mType metric.Type, id string) (val metric.Metric, ok bool) {
	ctx, cancel := context.WithTimeout(ctx, timeoutReqs)
	defer cancel()

	var mDelta sql.NullInt64
	var mValue sql.NullFloat64

	query := `SELECT delta, value FROM metrics WHERE id = $1 AND type = $2`
	row := s.db.QueryRowContext(ctx, query, id, mType.ToString())

	err := row.Scan(&mDelta, &mValue)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			s.log.Error(err)
		}

		return
	}

	val = metric.Metric{
		ID:    id,
		MType: mType,
	}

	if mDelta.Valid {
		val.Delta = new(int64)
		*val.Delta = mDelta.Int64
	}

	if mValue.Valid {
		val.Value = new(float64)
		*val.Value = mValue.Float64
	}

	return val, true
}

// GetAll - получить все метрики
func (s *Storage) GetAll(ctx context.Context) []metric.Metric {
	ctx, cancel := context.WithTimeout(ctx, timeoutReqs)
	defer cancel()

	query := `SELECT id, delta, value, type FROM metrics`
	rows, err := s.db.QueryContext(ctx, query)
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			s.log.Error(closeErr)
		}
	}()

	if err != nil {
		s.log.Error(err)

		return nil
	}

	var metrics []metric.Metric
	for rows.Next() {
		var mDelta sql.NullInt64
		var mValue sql.NullFloat64
		var id string
		var mType metric.Type

		err := rows.Scan(&id, &mDelta, &mValue, &mType)

		if err != nil {
			s.log.Error(err)

			return nil
		}

		m := metric.Metric{
			ID:    id,
			MType: mType,
		}

		if mDelta.Valid {
			m.Delta = new(int64)
			*m.Delta = mDelta.Int64
		}

		if mValue.Valid {
			m.Value = new(float64)
			*m.Value = mValue.Float64
		}

		metrics = append(metrics, m)
	}

	if err := rows.Err(); err != nil {
		s.log.Error(err)

		return nil
	}

	return metrics
}

func (s *Storage) setHandler(ctx context.Context, mType metric.Type, val metric.Metric) func() error {
	return func() error {
		ctxWithTTL, cancel := context.WithTimeout(ctx, timeoutReqs)
		defer cancel()

		query := `
INSERT INTO metrics (id, type, delta, value)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id, type)
DO UPDATE SET delta = $3, value = $4;
`
		var mDelta sql.NullInt64
		var mValue sql.NullFloat64

		if val.Delta != nil {
			mDelta.Valid = true
			mDelta.Int64 = *val.Delta
		}

		if val.Value != nil {
			mValue.Valid = true
			mValue.Float64 = *val.Value
		}

		_, err := s.db.ExecContext(
			ctxWithTTL,
			query,
			val.ID,
			mType.ToString(),
			mDelta,
			mValue,
		)

		return err
	}
}

func (s *Storage) batchSetHandler(ctx context.Context, mSlice []metric.Metric) func() error {
	return func() error {
		ctxWithTTL, cancel := context.WithTimeout(ctx, timeoutReqs)
		defer cancel()

		tx, err := s.db.BeginTx(ctxWithTTL, nil)
		if err != nil {
			return err
		}

		query := `
INSERT INTO metrics (id, type, delta, value)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id, type)
DO UPDATE SET delta = $3, value = $4;	
`

		stmt, err := tx.PrepareContext(ctxWithTTL, query)
		if err != nil {
			return err
		}

		for _, m := range mSlice {
			var mDelta sql.NullInt64
			var mValue sql.NullFloat64

			if m.Delta != nil {
				mDelta.Valid = true
				mDelta.Int64 = *m.Delta
			}

			if m.Value != nil {
				mValue.Valid = true
				mValue.Float64 = *m.Value
			}

			if _, err := stmt.ExecContext(ctxWithTTL, m.ID, m.MType.ToString(), mDelta, mValue); err != nil {
				return err
			}
		}

		return tx.Commit()
	}
}
