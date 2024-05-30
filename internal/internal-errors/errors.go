package internalerrors

import "errors"

var ErrInternal error = errors.New("Server internal error")
