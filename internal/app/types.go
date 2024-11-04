package app

import "wordofwisdom/internal/handler"

type Router interface {
	GetRoute(handlerID byte) (handler.HandlerFunc, bool)
}

type QuoteRepository interface {
	GetQuote() (string, error)
}

type PoWAlgorithm interface {
	GenerateChallenge() []byte
	ValidateSolution(challenge []byte, nonce []byte) bool
	Solve(challenge []byte) []byte
}
