package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"goflow/backend/internal/repository"
)

func mapErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.ErrNotFound
	}
	return err
}

func normPage(p repository.Page) (limit, offset int) {
	limit = p.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	offset = p.Offset
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func normMessageLimit(n int) int {
	if n <= 0 {
		return 50
	}
	if n > 200 {
		return 200
	}
	return n
}
