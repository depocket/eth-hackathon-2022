package utils

import "errors"

var ErrRecordNotFound = errors.New("record not found")

func IsNotFound(err error) bool {
	return errors.Is(err, ErrRecordNotFound)
}
