package utils

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

func HandelPgError(err error, notFoundErr error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return notFoundErr
	}

	return err
}
