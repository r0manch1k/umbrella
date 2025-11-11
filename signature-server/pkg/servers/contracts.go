package servers

import "context"

type Server interface {
	Start()
	Notify() <-chan error
	Shutdown() error
	Context() context.Context
}
