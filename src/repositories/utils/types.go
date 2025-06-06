package utils

import (
	"errors"
	"fmt"
)

var NilTxError = errors.New("nil transaction")

func NotFoundErrorF(resourceName string) error {
	return fmt.Errorf("resource: %s not found", resourceName)
}

var RowLengthShouldBeAtMost1Error = errors.New("length of rows should be at most 1")
