package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"sync"
	"time"

	"wordofwisdom/internal/service/pow"
	"wordofwisdom/pkg/config"
)

type PoWServer struct {
	ctx       context.Context
	cancel    context.CancelFunc
	config    *config.ServerConfig
	repo      QuoteRepository
	powAlgo   PoWAlgorithm
	connQueue chan net.Conn
	wg        sync.WaitGroup
	router    Router
}

func NewPoWServer(
	ctx context.Context,
	repo QuoteRepository,
	router Router,
) *PoWServer {
	cfg, err := config.ServerParse()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	SetLogger(cfg.Server.LogLevel)

	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)

	server := &PoWServer{
		ctx:       ctx,
		cancel:    cancel,
		config:    cfg,
		powAlgo:   pow.NewProofOfWork(cfg.POW.Complexity),
		repo:      repo,
		connQueue: make(chan net.Conn, 100),
		router:    router,
	}

	return server
}

func (s *PoWServer) Start() error {
	listener, err := net.Listen("tcp", s.config.Server.Address)
	if err != nil {
		slog.Error("failed to start listening", "error", err)
		return fmt.Errorf("failed to listen on %s: %w", s.config.Server.Address, err)
	}
	defer listener.Close()
	slog.Info("server started", "addr", s.config.Server.Address)

	// run workers
	for i := 0; i < s.config.Server.WorkerCount; i++ {
		s.wg.Add(1)
		go s.worker(i)
		slog.Info("worker started", "worker_id", i)
	}

	// accept incoming connections
	for {
		select {
		case <-s.ctx.Done():
			slog.Info("server context cancelled, stopping accepting new connections")
			return nil
		default:
			var conn net.Conn
			if conn, err = listener.Accept(); err != nil {
				slog.Warn("failed to accept connection", "error", err)
				continue
			}
			slog.Info("connection accepted", "remote_addr", conn.RemoteAddr().String())
			s.connQueue <- conn
		}
	}
}

func (s *PoWServer) worker(workerID int) {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case conn, ok := <-s.connQueue:
			if !ok {
				slog.Debug("connection queue closed, stopping worker", "worker_id", workerID)
				return
			}
			slog.Debug("worker processing connection", "worker_id", workerID)
			s.handleConnection(conn)
		}
	}
}

func (s *PoWServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	done := make(chan struct{})
	defer close(done)

	// send PoW challenge to client
	challenge := s.powAlgo.GenerateChallenge()

	// create a buffer for challenge + difficulty
	message := append(challenge, s.config.POW.Complexity)
	_, err := conn.Write(message)
	if err != nil {
		slog.Error("failed to send challenge", "error", err)
		return
	}
	slog.Debug("challenge and difficulty sent", "challenge", challenge, "difficulty", s.config.POW.Complexity)

	// Step 2: Read client solution (binary) and handler_id (as binary data)
	go func() {
		solution := make([]byte, 8) // Expecting 8 bytes for the solution
		_, err = conn.Read(solution)
		if err != nil {
			slog.Error("failed to read solution", "error", err)
			return
		}
		slog.Debug("solution received", "solution", solution)

		// Read handler_id as binary (assuming it's a single byte or a small number)
		handlerIDBuf := make([]byte, 1)
		_, err = conn.Read(handlerIDBuf)
		if err != nil {
			slog.Error("failed to read handler ID", "error", err)
			return
		}
		handlerID := handlerIDBuf[0]
		slog.Debug("handler ID received", "handler_id", handlerID)

		// Step 3: Validate PoW solution
		if s.powAlgo.ValidateSolution(challenge, solution) {
			slog.Debug("PoW solution validated", "solution", solution)

			// Step 4: Route to the appropriate handler using byte key
			handler, exists := s.router.GetRoute(handlerID)

			if !exists {
				if _, err = conn.Write([]byte("Handler not found\n")); err != nil {
					slog.Error("failed to write to connection", "error", err)
					return
				}
				slog.Warn("handler not found", "handler_id", handlerID)
				return
			}

			// Call the handler
			handler(conn, s.repo)
			slog.Debug("handler executed", "handler_id", handlerID)
		} else {
			if _, err = conn.Write([]byte("Invalid PoW solution\n")); err != nil {
				slog.Error("failed to write to connection", "error", err)
				return
			}
			slog.Warn("invalid PoW solution", "solution", solution)
		}

		done <- struct{}{}
	}()

	// Step 6: handle timeout
	select {
	case <-done:
		slog.Debug("connection handled successfully")
	case <-time.After(s.config.POW.Timeout):
		slog.Warn("timeout waiting for PoW solution")
		if _, err = conn.Write([]byte("Timeout waiting for PoW solution\n")); err != nil {
			slog.Error("failed to write timeout message", "error", err)
		}
	}
}

func (s *PoWServer) Shutdown() error {
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), s.config.Server.ShutdownTimeout)
	defer cancel()

	// Cancel the server's context to stop workers and other operations
	s.cancel()

	// Close connection queue
	close(s.connQueue)

	// Wait for workers to finish processing connections
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	// Wait for workers or context timeout
	select {
	case <-done:
		slog.Info("all workers have finished")
	case <-ctx.Done():
		slog.Warn("context timeout exceeded during shutdown")
	}

	return nil
}
