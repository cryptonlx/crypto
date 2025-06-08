package utils

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

var NilTxError = errors.New("nil transaction")
var UniqueViolationError = errors.New("unique_violation")

func NotFoundErrorF(resourceName string) error {
	return fmt.Errorf("resource: %s not found", resourceName)
}
func ConstraintViolationErrorF(constraintName string) error {
	if constraintName == "wallets_balance_check" {
		return fmt.Errorf("insufficient_funds")
	}
	return fmt.Errorf("constraint violation: %s", constraintName)
}

var RowLengthShouldBeAtMost1Error = errors.New("length of rows should be at most 1")

func ToError(err *pgconn.PgError) error {
	if err == nil {
		return nil
	}

	switch err.Code {
	case "23505":
		return UniqueViolationError
	case "23514":
		return ConstraintViolationErrorF(err.ConstraintName)
	}

	return err
}
