package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

const timeoutReqs = 10 * time.Second

type Storage struct {
	db  *sql.DB
	log *logger.Logger
}

func New(db *sql.DB, log *logger.Logger) *Storage {
	return &Storage{db: db, log: log}
}

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

func (s *Storage) Set(mType metric.Type, val metric.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutReqs)
	defer cancel()

	query := `
INSERT INTO metrics (id, type, delta, value)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id, type)
DO UPDATE SET delta = $3, value = $4;	
`
	mDelta := sql.NullInt64{}
	mValue := sql.NullFloat64{}

	if val.Delta != nil {
		mDelta.Valid = true
		mDelta.Int64 = *val.Delta
	}

	if val.Value != nil {
		mValue.Valid = true
		mValue.Float64 = *val.Value
	}

	_, err := s.db.ExecContext(
		ctx,
		query,
		val.ID,
		mType.ToString(),
		mDelta,
		mValue,
	)

	if err != nil {
		s.log.Error(err)
	}
}

func (s *Storage) Get(mType metric.Type, id string) (val metric.Metric, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutReqs)
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

func (s *Storage) GetAllByType(mType metric.Type) (values map[string]metric.Metric, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutReqs)
	defer cancel()

	query := `SELECT id, delta, value FROM metrics WHERE type = $1`
	rows, err := s.db.QueryContext(ctx, query, mType.ToString())
	defer func() {
		err := rows.Close()
		if err != nil {
			s.log.Error(err)
		}
	}()

	if err != nil {
		s.log.Error(err)

		return
	}

	tmpValues := make(map[string]metric.Metric)

	for rows.Next() {
		var mDelta sql.NullInt64
		var mValue sql.NullFloat64
		var id string

		err := rows.Scan(&id, &mDelta, &mValue)

		if err != nil {
			s.log.Error(err)

			return
		}

		val := metric.Metric{
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

		tmpValues[id] = val
	}

	if err := rows.Err(); err != nil {
		s.log.Error(err)

		return
	}

	values = tmpValues

	return values, true
}
