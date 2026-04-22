package idgen

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GenerateSequentialID(ctx context.Context, db *pgxpool.Pool, table, column, prefix string) (string, error) {
	var lastID string

	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY created_at DESC LIMIT 1", column, table)

	err := db.QueryRow(ctx, query).Scan(&lastID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Sprintf("%s-%04d", prefix, 1), nil
		}
		return "", err
	}

	var num int
	fmt.Sscanf(lastID, prefix+"-%04d", &num)
	return fmt.Sprintf("%s-%04d", prefix, num+1), nil
}

func GenerateDailyID(ctx context.Context, db *pgxpool.Pool, table, column, prefix string) (string, error) {
	today := time.Now().UTC()
	dateStr := today.Format("2006/01/02")

	var count int
	query := fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE DATE(created_at) = CURRENT_DATE`,
		table,
	)

	err := db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		count = 0
	}

	return fmt.Sprintf("%s/%s/%03d", prefix, dateStr, count+1), nil
}
