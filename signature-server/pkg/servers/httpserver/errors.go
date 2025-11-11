package httpserver

import "errors"

var ErrShutdownTimeout = errors.New("shutdown timeout exceeded")
