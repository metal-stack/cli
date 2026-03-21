package helpers

import (
	"errors"

	"connectrpc.com/connect"
)

func IsAlreadyExists(err error) bool {
	if connectErr, ok := errors.AsType[*connect.Error](err); ok && connectErr.Code() == connect.CodeAlreadyExists {
		return true
	}

	return false
}
