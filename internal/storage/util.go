package storage

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func DurationToPGInterval(duration time.Duration) pgtype.Interval {
	totalMicroseconds := duration.Microseconds()
	days := totalMicroseconds / (24 * 60 * 60 * 1_000_000)
	microsecondsRemaining := totalMicroseconds % (24 * 60 * 60 * 1_000_000)

	return pgtype.Interval{
		Months:       0, // Assuming no months component
		Days:         int32(days),
		Microseconds: microsecondsRemaining,
		Valid:        true,
	}
}

func TimeToPgxTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}
