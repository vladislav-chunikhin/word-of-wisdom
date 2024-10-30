package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"sync"
	"time"

	"lib-go/pkg/config"
)

const (
	defaultNetworkProtocol = "tcp"
)

type Server struct {
	cfg      *config.Config
	listener net.Listener

	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func NewServer() *Server {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	setLogger(cfg.Server.LogLevel)

	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Run(ctx context.Context, connectionHandler func(conn net.Conn) error) (err error) {
	if connectionHandler == nil {
		return fmt.Errorf("connection handler is not defined")
	}

	ctx, s.cancel = context.WithCancel(ctx)
	defer s.cancel()

	lc := net.ListenConfig{
		KeepAlive: s.cfg.Server.KeepAlive,
	}

	s.listener, err = lc.Listen(ctx, defaultNetworkProtocol, s.cfg.Server.Address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	slog.Debug(fmt.Sprintf("server started on port %s", s.listener.Addr().String()))

	s.wg.Add(1)
	go s.serve(ctx, connectionHandler)
	s.wg.Wait()

	slog.Debug("server stopped")

	return nil
}

func (s *Server) Config() *config.Config {
	return s.cfg
}

func (s *Server) serve(ctx context.Context, connectionHandler func(conn net.Conn) error) {
	defer s.wg.Done()

	go func() {
		<-ctx.Done()
		err := s.listener.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			slog.Error(fmt.Sprintf("failed to close listener: %v", err))
		}
	}()

	for {
		conn, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			slog.Debug("listener closed")
			return
		} else if err != nil {
			slog.Error(fmt.Sprintf("failed to accept connection: %v", err))
			continue
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.handleConnection(conn, connectionHandler)
		}()
	}
}

func (s *Server) handleConnection(conn net.Conn, connectionHandler func(conn net.Conn) error) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("failed to close connection: %v", err))
		}
	}(conn)

	if err := conn.SetDeadline(time.Now().Add(s.cfg.Server.Deadline)); err != nil {
		slog.Error(fmt.Sprintf("failed to set deadline: %v", err))
		return
	}

	if err := connectionHandler(conn); err != nil {
		slog.Error(fmt.Sprintf("failed to handle connection: %v", err))
	}
}
