package repository

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"math/rand"
	"strings"
)

//go:embed file/quotes.txt
var quotes []byte

type Quote struct {
	quotes []string
}

func NewQuote() *Quote {
	qr := new(Quote)

	reader := bytes.NewReader(quotes)
	s := bufio.NewScanner(reader)

	for s.Scan() {
		if q := strings.TrimSpace(s.Text()); q != "" {
			qr.quotes = append(qr.quotes, q)
		}
	}

	return qr
}

func (q *Quote) GetQuote() (string, error) {
	if len(q.quotes) == 0 {
		return "", errors.New("no quotes available")
	}

	return q.quotes[rand.Intn(len(q.quotes))], nil
}
