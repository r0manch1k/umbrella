package httpserver

import (
	"net"
	"time"
)

type Option func(*Server)

func Address(host, port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort(host, port)
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
