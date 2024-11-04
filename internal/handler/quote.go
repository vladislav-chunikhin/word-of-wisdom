package handler

import (
	"errors"
	"log/slog"
	"net"
)

func HandleQuote(conn net.Conn, repo QuoteRepository) {
	quote, err := repo.GetQuote()
	if err != nil {
		_, writeErr := conn.Write([]byte("failed to get quote: " + err.Error() + "\n"))
		if writeErr != nil {
			err = errors.Join(err, writeErr)
			slog.Error("error writing to connection", "error", err)
		}
		return
	}

	if _, err = conn.Write([]byte("here is your quote: '" + quote + "\n")); err != nil {
		slog.Error("error writing to connection", "error", err)
	}
}
