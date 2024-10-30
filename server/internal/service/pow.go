package service

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"net"

	"lib-go/pkg/pow"
)

var (
	ErrInvalidSolution  = errors.New("invalid solution")
	ErrDiffChallenge    = errors.New("got different challenge")
	ErrFailedToGetQuote = errors.New("failed to get quote")
)

type QuotesRepository interface {
	GetQuote() (string, error)
}

type PowProvider interface {
	GenerateChallenge() []byte
	CheckSolution(solution pow.PowSolution) bool
}

type POW struct {
	quotesRepo  QuotesRepository
	powProvider PowProvider
}

func NewPOW(quotesRepo QuotesRepository, powProvider PowProvider) *POW {
	return &POW{quotesRepo: quotesRepo, powProvider: powProvider}
}

func (p *POW) Handle(conn net.Conn) error {
	clientAddress := conn.RemoteAddr().String()
	slog.Debug("connected new client", "client_address", clientAddress)

	originalChallenge := p.powProvider.GenerateChallenge()
	if _, err := conn.Write(originalChallenge); err != nil {
		return err
	}

	buffer := make([]byte, 16)
	if _, err := io.ReadFull(conn, buffer); err != nil {
		return err
	}
	nonce := binary.BigEndian.Uint64(buffer[:8])
	challenge := buffer[8:]
	if !bytes.Equal(originalChallenge[:8], challenge) {
		return ErrDiffChallenge
	}

	solution := pow.NewSolution(challenge, nonce)
	result := p.powProvider.CheckSolution(solution)
	if !result {
		return ErrInvalidSolution
	}

	quote, err := p.quotesRepo.GetQuote()
	if err != nil {
		slog.Error("failed to get quote", "client_address", clientAddress, "error", err)
		return ErrFailedToGetQuote
	}

	if _, err = conn.Write([]byte(quote)); err != nil {
		return err
	}

	slog.Debug("sent quote", "client_address", clientAddress)
	return nil
}
