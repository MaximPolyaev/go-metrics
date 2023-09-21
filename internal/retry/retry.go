package retry

import (
	"errors"
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func Retry(attemptHandler func() error, log *logger.Logger) {
	var prevT *time.Duration

	for {
		if prevT != nil {
			time.Sleep(*prevT)
		}

		if err := attemptHandler(); err != nil {
			log.Error(err)

			if isRetryReq(err) {
				prevT = getNextAttemptTime(prevT)
			} else {
				prevT = nil
			}

			if prevT != nil {
				continue
			}
		}

		break
	}
}

func getNextAttemptTime(prevT *time.Duration) *time.Duration {
	if prevT == nil {
		t := time.Second
		return &t
	}

	switch *prevT {
	case time.Second:
		t := 3 * time.Second
		return &t
	case 3 * time.Second:
		t := 5 * time.Second
		return &t
	default:
		return nil
	}
}

func isRetryReq(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ConnectionException
}
