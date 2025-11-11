package httpserver

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/r0manch1k/umbrella/signature-server/pkg/logger"
	"github.com/r0manch1k/umbrella/signature-server/pkg/servers"
	"github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

var _ servers.Server = (*Server)(nil)

// Server — обёртка над fasthttp.Server.
type Server struct {
	ctx context.Context
	eg  *errgroup.Group

	Server *fasthttp.Server
	notify chan error

	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration

	logger logger.Interface
}

// New — создаёт новый экземпляр HTTPServer.
func New(l logger.Interface, handler fasthttp.RequestHandler, serverName string, opts ...Option) *Server {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1)

	s := &Server{
		ctx:             ctx,
		eg:              group,
		notify:          make(chan error, 1),
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
		logger:          l,
	}

	wrappedHandler := accessLogger(l, handler)

	for _, opt := range opts {
		opt(s)
	}

	s.Server = &fasthttp.Server{
		Handler:      wrappedHandler,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		Name:         serverName,
	}

	return s
}

func accessLogger(l logger.Interface, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		next(ctx)

		var accessLog strings.Builder

		accessLog.WriteString(string(ctx.Request.Header.Method()))
		accessLog.WriteString(" ")
		accessLog.WriteString(string(ctx.Request.URI().PathOriginal()))
		accessLog.WriteString(" - ")
		accessLog.WriteString(strconv.Itoa(ctx.Response.StatusCode()))
		accessLog.WriteString(" - ")
		accessLog.WriteString(ctx.RemoteIP().String())
		accessLog.WriteString(" ")
		accessLog.WriteString(strconv.Itoa(len(ctx.Request.Body())))

		l.Info(accessLog.String())
	}
}

// Start — запускает сервер асинхронно.
func (s *Server) Start() {
	s.eg.Go(func() error {
		err := s.Server.ListenAndServe(s.address)
		if err != nil {
			s.notify <- err

			close(s.notify)

			return err
		}

		return nil
	})

	s.logger.Info("HTTP server - Started at %s", s.address)
}

// Notify — возвращает канал с ошибками сервера.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown — корректно завершает работу сервера.
func (s *Server) Shutdown() error {
	var shutdownErrors []error

	s.logger.Info("HTTP server - Shutting down...")

	ctx, cancel := context.WithTimeout(s.ctx, s.shutdownTimeout)
	defer cancel()

	shutdownDone := make(chan struct{})

	go func() {
		if err := s.Server.Shutdown(); err != nil {
			shutdownErrors = append(shutdownErrors, err)
		}

		close(shutdownDone)
	}()

	select {
	case <-ctx.Done():
		shutdownErrors = append(shutdownErrors, ErrShutdownTimeout)
	case <-shutdownDone:
	}

	err := s.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		shutdownErrors = append(shutdownErrors, err)
	}

	s.logger.Info("HTTP server - Shutdown complete")

	return errors.Join(shutdownErrors...)
}

// Context — возвращает контекст сервера.
func (s *Server) Context() context.Context {
	return s.ctx
}
